/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kcoidc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"stash.kopano.io/kgol/oidc-go"

	"stash.kopano.io/kc/libkcoidc/internal/version"
)

// A Provider is a representation of an OpenID Connect Provider (OP).
type Provider struct {
	mutex sync.RWMutex

	initialized bool
	provider    *oidc.Provider
	ready       chan struct{}

	httpClient *http.Client

	logger Logger
	debug  bool

	definition *oidc.ProviderDefinition
}

var emptyProviderDefintion = &oidc.ProviderDefinition{}

// NewProvider creates a new Provider with the provider HTTP client. If no client
// is provided, http.DefaultClient will be used.
func NewProvider(client *http.Client, logger Logger, debug bool) (*Provider, error) {
	if client == nil {
		client = http.DefaultClient
	}

	if logger == nil {
		logger = DefaultLogger
	}

	p := &Provider{
		httpClient: client,

		logger: logger,
		debug:  debug,
	}
	return p, nil
}

// Version returns the runtime version string of this module.
func (p *Provider) Version() string {
	return version.Version
}

// BuildDate returns the build data string of this module.
func (p *Provider) BuildDate() string {
	return version.BuildDate
}

// Initialize initializes the associated Provider with the provided issuer.
func (p *Provider) Initialize(ctx context.Context, issuer *url.URL) error {
	var err error
	if issuer.Host == "" || issuer.Scheme != "https" {
		return ErrStatusInvalidIss
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.initialized {
		return ErrStatusAlreadyInitialized
	}

	ready := make(chan struct{})
	p.ready = ready

	updates := make(chan *oidc.ProviderDefinition)
	config := &oidc.ProviderConfig{
		HTTPClient: p.httpClient,
		Logger:     p.logger,
	}
	provider, err := oidc.NewProvider(issuer, config)
	if err != nil {
		if p.logger != nil {
			p.logger.Printf("kcoidc initialize failed to created provider: %v", err)
		}
		return ErrStatusInvalidIss
	}
	err = provider.Initialize(ctx, updates, nil)
	if err != nil {
		if p.logger != nil {
			p.logger.Printf("kcoidc initialize failed with error: %v", err)
		}
		return ErrStatusUnknown
	}

	p.provider = provider
	p.definition = emptyProviderDefintion
	p.initialized = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				if update == nil {
					return
				}
				p.mutex.Lock()
				d := p.definition
				p.definition = update
				p.mutex.Unlock()
				if d == emptyProviderDefintion {
					close(ready)
				}
			}
		}
	}()
	return nil
}

// Uninitialize uninitializes the associated Provider.
func (p *Provider) Uninitialize() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.initialized {
		return ErrStatusNotInitialized
	}

	err := p.provider.Shutdown()
	if p.logger != nil {
		p.logger.Printf("kcoidc provider shutdown with error: %v", err)
	}
	p.initialized = false
	p.provider = nil

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

// ValidateTokenString validates the provided token string value with the keys
// of the accociated Provider and returns the authenticated users ID as found in
// the claims, the standard claims and all extra claims.
func (p *Provider) ValidateTokenString(ctx context.Context, tokenString string) (string, *jwt.StandardClaims, *ExtraClaimsWithType, error) {
	p.mutex.RLock()
	ddoc := p.definition.WellKnown
	jwks := p.definition.JWKS
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
		for _, alg := range ddoc.IDTokenSigningAlgValuesSupported {
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
			p.logger.Printf("kcoidc validate token key: %#v (%v)\n", key.Key, kid)
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
	p.mutex.RLock()
	ddoc := p.definition.WellKnown
	p.mutex.RUnlock()
	if ddoc == nil {
		return nil, ErrStatusNotInitialized
	}

	var contentTypeJSONOnly = []string{"application/json"}
	userinfo := make(map[string]interface{})

	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", tokenString)},
	}

	return userinfo, fetchJSON(ctx, p.httpClient, ddoc.UserInfoEndpoint, headers, contentTypeJSONOnly, &userinfo)
}
