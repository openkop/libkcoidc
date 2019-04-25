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

package main

import (
	"C"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/http2"

	"stash.kopano.io/kc/libkcoidc"
)

// Global library state. This also means that this library can only use a single
// OIDC Provider at the same time as the issuer is directly bound to the global
// library state.
var (
	mutex                    sync.RWMutex
	client                   *http.Client
	transport                *http.Transport
	initializedContext       context.Context
	initializedContextCancel context.CancelFunc
	initializedLogger        *log.Logger
	debug                    bool
	provider                 *kcoidc.Provider
)

func init() {
	if os.Getenv("KCOIDC_DEBUG") != "" {
		debug = true
		fmt.Println("kcoidc-c debug enabled")
		initializedLogger = log.New(os.Stdout, "kcoidc-c debug ", 0)
	}

	// TODO(longsleep): Add HTTP client env vars same as kcc-go/http.go.

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: tls.NewLRUClientSessionCache(0),
		},
	}

	// Enable HTTP2 support.
	err := http2.ConfigureTransport(transport)
	if err != nil {
		panic(err)
	}

	client = &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}

	// Setup transport defaults.
	InsecureSkipVerify(false)
}

// Initialize initializes the global library state with the provided issuer.
func Initialize(ctx context.Context, iss string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if provider != nil {
		return kcoidc.ErrStatusAlreadyInitialized
	}

	issURL, err := url.Parse(iss)
	if err != nil {
		if debug {
			fmt.Printf("kcoidc-c initialize failed with invalid iss value: %v\n", err)
		}
		return kcoidc.ErrStatusInvalidIss
	}

	var p *kcoidc.Provider
	if initializedLogger == nil {
		p, err = kcoidc.NewProvider(client, nil, debug)
	} else {
		p, err = kcoidc.NewProvider(client, initializedLogger, debug)
	}
	if err != nil {
		if debug {
			fmt.Printf("kcoidc-c initialize failed: %v\n", err)
		}
		return err
	}

	err = p.Initialize(ctx, issURL)
	if err != nil {
		if debug {
			fmt.Printf("kcoidc-c initialize failed: %v\n", err)
		}
		return err
	}

	provider = p
	initializedContext, initializedContextCancel = context.WithCancel(ctx)
	if debug {
		fmt.Printf("kcoidc-c initialize success: %v\n", iss)
	}
	return nil
}

// Uninitialize uninitializes the global library state.
func Uninitialize() error {
	mutex.Lock()
	defer mutex.Unlock()

	if provider == nil {
		return nil
	}

	if debug {
		fmt.Println("kcoidc-c uninitialize")
	}

	err := provider.Uninitialize()
	if err != nil {
		return err
	}

	initializedContextCancel()
	initializedContext = nil
	initializedContextCancel = nil

	provider = nil
	if debug {
		fmt.Println("kcoidc-c uninitialize success")
	}
	return nil
}

// InsecureSkipVerify sets up the libraries HTTP transport according to the
// provided parametters.
func InsecureSkipVerify(insecureSkipVerify bool) error {
	mutex.RLock()
	defer mutex.RUnlock()

	if provider != nil {
		return kcoidc.ErrStatusAlreadyInitialized
	}

	if insecureSkipVerify != transport.TLSClientConfig.InsecureSkipVerify {
		if insecureSkipVerify {
			transport.TLSClientConfig.InsecureSkipVerify = true
			if debug {
				fmt.Println("kcoidc-c TLS verification is now disabled - this is insecure")
			}
		} else {
			transport.TLSClientConfig.InsecureSkipVerify = false
			if debug {
				fmt.Println("kcoidc-c TLS verification is now enabled")
			}
		}
	}

	return nil
}

// WaitUntilReady blocks until the initialization is ready or timeout.
func WaitUntilReady(timeout time.Duration) error {
	mutex.RLock()
	p := provider
	ctx := initializedContext
	mutex.RUnlock()

	var err error
	if debug {
		fmt.Println("kcoidc-c waiting until ready")
		defer func() {
			fmt.Printf("kcoidc-c finished waiting until ready: %v\n", err)
		}()
	}

	if p == nil {
		err = kcoidc.ErrStatusNotInitialized
	} else {
		err = p.WaitUntilReady(ctx, timeout)
	}

	return err
}

// ValidateTokenString validates the provided token string value and returns
// the authenticated users ID as found the claims the standard claims and all
// extra claims. Error will be set when the validation failed.
func ValidateTokenString(tokenString string) (string, *jwt.StandardClaims, *kcoidc.ExtraClaimsWithType, error) {
	mutex.RLock()
	p := provider
	ctx := initializedContext
	mutex.RUnlock()

	if debug {
		fmt.Printf("kcoidc-c validate token string: %s\n", tokenString)
	}
	if p == nil {
		return "", nil, nil, kcoidc.ErrStatusNotInitialized
	}

	authenticatedUserID, standardClaims, extraClaims, err := p.ValidateTokenString(ctx, tokenString)
	if err != nil && debug {
		fmt.Printf("kcoidc-c validate token resulted in validation failure: %s\n", err)
	}
	return authenticatedUserID, standardClaims, extraClaims, err
}

// ValidateTokenStringAndRequireClaim validates the provided token string value
//  and returns the authenticated users ID as found the claims the standard
// claims and all extra claims. In addition, the token must have authenticated
// the provided requiredScope. Error will be set when the validation failed or
// the required scope is not authenticated.
func ValidateTokenStringAndRequireClaim(tokenString string, requiredScope string) (string, *jwt.StandardClaims, *kcoidc.ExtraClaimsWithType, error) {
	authenticatedUserID, standardClaims, extraClaims, err := ValidateTokenString(tokenString)
	if err != nil {
		return authenticatedUserID, standardClaims, extraClaims, err
	}

	err = kcoidc.RequireScopesInClaims(extraClaims, []string{requiredScope})
	if err != nil && debug {
		fmt.Printf("kcoidc-c validate token and require claims result in scope require failure: %s\n", err)
	}

	return authenticatedUserID, standardClaims, extraClaims, err
}

// FetchUserinfoWithAccesstokenString fetches the available user info for the
// provided access token and returns it as a string map of values.
func FetchUserinfoWithAccesstokenString(tokenString string) (map[string]interface{}, error) {
	mutex.RLock()
	p := provider
	ctx := initializedContext
	mutex.RUnlock()

	if p == nil {
		return nil, kcoidc.ErrStatusNotInitialized
	}

	userinfo, err := p.FetchUserinfoWithAccesstokenString(ctx, tokenString)
	if err != nil && debug {
		fmt.Printf("kcoidc-c fetch userinfo failure: %s\n", err)
	}
	return userinfo, err
}

func main() {}
