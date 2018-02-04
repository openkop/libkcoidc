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
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"stash.kopano.io/kc/libkcoidc"
)

func run(issString, tokenString string) error {
	ctx := context.Background()

	// Initialize with insecure operations allowed.
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	provider, err := kcoidc.NewProvider(client, nil, false)
	if err != nil {
		fmt.Printf("> Error: failed to create provider: %v\n", err)
		return err
	}
	// Initialize with issuer identifier.
	issURL, err := url.Parse(issString)
	if err != nil {
		fmt.Printf("> Error: failed to parse issuer: %v\n", err)
		return err
	}
	err = provider.Initialize(ctx, issURL)
	if err != nil {
		fmt.Printf("> Error: initialize failed: %v\n", err)
		return err
	}
	// Wait until oidc validation becomes ready.
	err = provider.WaitUntilReady(ctx, 10*time.Second)
	if err != nil {
		fmt.Printf("> Error: failed to get ready in time: %v\n", err)
		return err
	}

	beginTime := time.Now()
	sub, standardClaims, extraClaims, err := provider.ValidateTokenString(ctx, tokenString)
	endTime := time.Now()
	duration := endTime.Sub(beginTime)

	validString := "valid"
	if err != nil {
		validString = "invalid"
	}

	if e := printResultOrError(err, "Result code"); e != nil {
		fmt.Printf("> Error: failed to validate token string: %v\n", e)
	}

	fmt.Printf("> Token subject : %s -> %s\n", sub, validString)
	fmt.Printf("> Time spent    : %fs\n", duration.Seconds())
	fmt.Printf("> Standard      : %v\n", standardClaims)
	fmt.Printf("> Extra         : %v\n", extraClaims)
	fmt.Printf("> Token type    : %d\n", extraClaims.KCTokenType())

	if err == nil && extraClaims.KCTokenType() == kcoidc.TokenTypeKCAccess {
		userinfo, userinfoErr := provider.FetchUserinfoWithAccesstokenString(ctx, tokenString)

		if e := printResultOrError(userinfoErr, "Userinfo   "); e != nil {
			fmt.Printf("> Error: failed to fetch userinfo: %v\n", e)
		} else {
			fmt.Printf("%v", userinfo)
		}
	}

	// Clean up as well.
	if e := provider.Uninitialize(); e != nil {
		fmt.Printf("> Error: failed to uninitialize: %v\n", err)
	}

	return nil
}

func printResultOrError(err error, msg string) error {
	switch e := err.(type) {
	case nil:
		fmt.Printf("> %s   : 0x0\n", msg)
	case kcoidc.ErrStatus:
		fmt.Printf("> %s   : 0x%x (%v)\n", msg, uint64(e), e)
	default:
		return err
	}

	return nil
}

func main() {
	var issString string
	var tokenString string

	if len(os.Args) > 1 {
		issString = os.Args[1]
	}
	if len(os.Args) > 2 {
		tokenString = os.Args[2]
	}

	err := run(issString, tokenString)
	if err != nil {
		os.Exit(-1)
	}
}
