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

	"stash.kopano.io/kc/libkcoidc"
)

// Global library state. This also means that this library can only use a single
// OIDC Provider at the same time as the issuer is directly bound to the global
// library state.
var (
	mutex                    sync.RWMutex
	client                   *http.Client
	initializedContext       context.Context
	initializedContextCancel context.CancelFunc
	logger                   *log.Logger
	debug                    bool
	provider                 *kcoidc.Provider
)

func init() {
	if os.Getenv("KCOIDC_DEBUG") != "" {
		debug = true
		fmt.Println("kcoidc-c debug enabled")
	}

	client = &http.Client{
		Timeout: 60 * time.Second,
	}

	logger = log.New(os.Stdout, "", 0)

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

	p, err := kcoidc.NewProvider(client, logger, debug)
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

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if insecureSkipVerify {
		// Only set this when we have something to change to allow Go to use
		// the internal HTTP2 connection logic otherwise.
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		if debug {
			fmt.Println("kcoidc-c TLS verification is now disabled - this is insecure")
		}
	}

	client.Transport = transport
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

	err = p.WaitUntilReady(ctx, timeout)
	return err
}

// ValidateTokenString validates the provided token string value and returns
// the subject as found in the claims.
func ValidateTokenString(tokenString string) (string, *jwt.StandardClaims, *kcoidc.ExtraClaimsWithType, error) {
	mutex.RLock()
	p := provider
	ctx := initializedContext
	mutex.RUnlock()

	if debug {
		fmt.Printf("kcoidc-c validate token string: %s\n", tokenString)
	}

	sub, standardClaims, extraClaims, err := p.ValidateTokenString(ctx, tokenString)
	if err != nil && debug {
		fmt.Printf("kcoid-c validate token resulted in validation failure: %s\n", err)
	}
	return sub, standardClaims, extraClaims, err
}

// FetchUserinfoWithAccesstokenString fetches the available user info for the
// provided access token and returns it as a string map of values.
func FetchUserinfoWithAccesstokenString(tokenString string) (map[string]interface{}, error) {
	mutex.RLock()
	p := provider
	ctx := initializedContext
	mutex.RUnlock()

	userinfo, err := p.FetchUserinfoWithAccesstokenString(ctx, tokenString)
	if err != nil && debug {
		fmt.Printf("kcoid-c fetch userinfo failure: %s\n", err)
	}
	return userinfo, err
}

func main() {}
