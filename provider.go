/*
 * Copyright 2018 Kopano and its licensors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3
 * or later, as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package kcoidc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// A Provider is a representation of an OpenID Connect Provider (OP).
type Provider struct {
	mutex sync.RWMutex

	initialized bool
	ready       chan struct{}
	started     chan error
	cancel      context.CancelFunc

	client *http.Client
	logger *log.Logger
	debug  bool

	iss       string
	discovery *oidcDiscoveryDocument
	jwks      *oidcJSONWebKeySet
}

// NewProvider creates a new Provider with the provider HTTP client. If no client
// is provided, http.DefaultClient will be used.
func NewProvider(client *http.Client, logger *log.Logger, debug bool) (*Provider, error) {
	if client == nil {
		client = http.DefaultClient
	}

	p := &Provider{
		client: client,
		logger: logger,
		debug:  debug,
	}
	return p, nil
}

// Initialize initializes the associated Provider with the provided issuer.
func (p *Provider) Initialize(ctx context.Context, iss *url.URL) error {
	var err error
	if iss.Host == "" || iss.Scheme == "" {
		return ErrStatusInvalidIss
	}

	p.mutex.Lock()
	if p.initialized {
		p.mutex.Unlock()
		return ErrStatusAlreadyInitialized
	}

	c, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.initialized = true

	p.iss = iss.String()

	started := make(chan error, 1)
	p.started = started
	go p.start(c, started)

	p.mutex.Unlock()

	err = <-started
	return err
}

// Uninitialize uninitializes the associated Provider.
func (p *Provider) Uninitialize() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.initialized {
		return ErrStatusNotInitialized
	}

	p.cancel()
	err := <-p.started

	p.cancel = nil
	p.started = nil
	p.iss = ""
	p.initialized = false
	p.ready = nil
	p.discovery = nil
	p.jwks = nil

	switch err {
	case context.Canceled:
		return nil
	}
	return err
}

// WaitUntilReady blocks until the associated Provider is ready or timeout.
func (p *Provider) WaitUntilReady(ctx context.Context, timeout time.Duration) error {
	p.mutex.RLock()
	if !p.initialized {
		p.mutex.RUnlock()
		return ErrStatusNotInitialized
	}
	ready := p.ready
	p.mutex.RUnlock()

	var err error
	select {
	case <-ready:
	case <-ctx.Done():
	case <-time.After(timeout):
		err = ErrStatusTimeout
	}

	return err
}

func (p *Provider) start(ctx context.Context, started chan error) {
	if p.debug && p.logger != nil {
		defer func() {
			p.logger.Println("kcoidc start has ended")
		}()
	}

	// Use started channel to signal caller that we are done.
	p.mutex.Lock()
	if !p.initialized || started != p.started {
		p.mutex.Unlock()
		started <- ErrStatusWrongInitialization
		return
	}

	// Create ready channel to keep ourselves running until success or another
	// signal makes us exit.
	ready := make(chan struct{})
	p.ready = ready
	p.mutex.Unlock()
	started <- nil

	for {
		retry := 60 * time.Second
		if p.debug && p.logger != nil {
			p.logger.Println("kcoidc running ...")
		}

		p.mutex.RLock()
		if p.initialized && started == p.started {
			iss := p.iss
			p.mutex.RUnlock()
			ddoc, err := newDiscoveryDocumentFromIssuer(ctx, p.client, iss)
			if err != nil {
				if p.logger != nil {
					p.logger.Printf("kcoid discovery error: %v\n", err)
				}
				retry = 5 * time.Second
			} else {
				jwks, err := newoidcJSONWebKeySetFromURL(ctx, p.client, ddoc.JWKSUri)
				if err != nil {
					if p.logger != nil {
						p.logger.Printf("kcoid discovery jwks error: %v\n", err)
					}
					retry = 5 * time.Second
				} else {
					p.mutex.Lock()
					if p.initialized && started == p.started {
						p.discovery = ddoc
						p.jwks = jwks
					}
					close(ready)
					p.mutex.Unlock()
				}
			}
		} else {
			p.mutex.RUnlock()
		}

		select {
		case <-ctx.Done():
			started <- ctx.Err()
			close(started)
			return
		case <-ready:
			close(started)
			return
		case <-time.After(retry):
			// We break for retries.
		}
	}
}

// ValidateTokenString validates the provided token string value with the keys
// of the accociated Provider and returns the authenticated users ID as found in
// the claims, the standard claims and all extra claims.
func (p *Provider) ValidateTokenString(ctx context.Context, tokenString string) (string, *jwt.StandardClaims, *ExtraClaimsWithType, error) {
	p.mutex.RLock()
	ddoc := p.discovery
	jwks := p.jwks
	p.mutex.RUnlock()
	if ddoc == nil || jwks == nil {
		return "", nil, nil, ErrStatusNotInitialized
	}

	claims := &ExtraClaimsWithType{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if p.debug && p.logger != nil {
			p.logger.Printf("kcoidc validate token header: %#v\n", token.Header)
		}

		supportedAlg := false
		for _, alg := range ddoc.SupportedSigningAlgValues {
			if token.Method.Alg() == alg {
				supportedAlg = true
				break
			}
		}
		if !supportedAlg {
			return nil, ErrStatusTokenUnexpectedSigningMethod
		}

		kid, _ := (token.Header["kid"].(string))
		keys := jwks.Key(kid)
		if keys == nil || len(keys) == 0 {
			return nil, ErrStatusTokenUnknownKey
		}

		key := keys[0]
		if p.debug && p.logger != nil {
			fmt.Printf("kcoidc validate token key: %#v (%v)\n", key.Key, kid)
		}

		return key.Key, nil
	})

	// Get standard claims.
	standardClaims, standardClaimsErr := SplitStandardClaimsFromMapClaims(claims)
	if err == nil {
		err = standardClaimsErr
	}
	if err == nil {
		err = standardClaims.Valid()
	}
	if err == nil && !token.Valid {
		// NOTE(longsleep): Can this actually happen?
		err = ErrStatusTokenValidationFailed
	}
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = ErrStatusTokenMalformed
			} else if ve.Errors&(jwt.ValidationErrorSignatureInvalid|jwt.ValidationErrorUnverifiable) != 0 {
				err = ErrStatusTokenInvalidSignature
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				err = ErrStatusTokenExpiredOrNotValidYet
			} else {
				err = ErrStatusTokenValidationFailed
			}
		}
	}

	// Get authenticated UserID
	authenticatedUserID, ok := AuthenticatedUserIDFromClaims(claims)
	if !ok {
		// NOTE(longsleep): Fallback to standard Subject if no extra information
		// is set in token. This can happen for older Konnect installations
		// which did not set this claim. Let's do this for compatibility.
		authenticatedUserID = standardClaims.Subject
	}

	return authenticatedUserID, standardClaims, claims, err
}

// FetchUserinfoWithAccesstokenString fetches the the userinfo result of the
// accociated provider for the provided access token string.
func (p *Provider) FetchUserinfoWithAccesstokenString(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	userinfo := make(map[string]interface{})

	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", tokenString)},
	}

	return userinfo, fetchJSON(ctx, p.client, p.discovery.UserinfoEndpoint, headers, &userinfo)
}
