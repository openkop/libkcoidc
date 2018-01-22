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
	"runtime"
	"sync"
	"time"

	"stash.kopano.io/kc/libkcoidc"
)

func benchValidateTokenS(provider *kcoidc.Provider, id int, count uint64, tokenString string) {
	fmt.Printf("> Info : thread %d started ...\n", id)

	var success uint64
	var failed uint64
	var i uint64
	for i = 0; i < count; i++ {
		_, err := provider.ValidateTokenString(tokenString)
		if err != nil {
			fmt.Printf("> Error: validation failed: %v\n", err)
			failed++
		} else {
			success++
		}
	}

	fmt.Printf("> Info : thread %d done:%d failed:%d ...\n", id, success, failed)
}

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

	concurrentThreadsSupported := runtime.NumCPU()
	var count uint64 = 100000
	var wg sync.WaitGroup

	// Wait until oidc validation becomes ready.
	err = provider.WaitUntilReady(10 * time.Second)
	if err != nil {
		fmt.Printf("> Error: failed to get ready in time: %v\n", err)
		return err
	}

	fmt.Printf("> Info : using %d threads with %d runs per thread\n", concurrentThreadsSupported, count)
	beginTime := time.Now()
	for id := 1; id <= concurrentThreadsSupported; id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			benchValidateTokenS(provider, id, count, tokenString)
		}(id)
	}
	wg.Wait()
	endTime := time.Now()
	duration := endTime.Sub(beginTime)
	rate := float64(count*uint64(concurrentThreadsSupported)) / duration.Seconds()
	fmt.Printf("> Time : %fs\n", duration.Seconds())
	fmt.Printf("> Rate : %f ops\n", rate)

	// Clean up as well.
	if e := provider.Uninitialize(); e != nil {
		fmt.Printf("> Error: failed to uninitialize: %v\n", err)
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
