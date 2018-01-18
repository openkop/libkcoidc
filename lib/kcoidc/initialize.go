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
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gopkg.in/square/go-jose.v2"
)

var initialization *initializationData

type initializationData struct {
	sync.RWMutex
	quit chan struct{}

	initialized bool
	ready       chan struct{}
	started     chan error
	cancel      context.CancelFunc

	client *http.Client

	iss       string
	discovery *oidcDiscoveryDocument
	jwks      *jose.JSONWebKeySet
}

func init() {
	initialization = &initializationData{
		quit: make(chan struct{}),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Initialize initializes the global library state with the provided issuer.
func Initialize(ctx context.Context, iss string) error {
	if debugEnabled {
		fmt.Printf("kcoidc initialize: %v\n", iss)
	}

	initialization.Lock()
	if initialization.initialized {
		initialization.Unlock()
		return KCOIDCErrAlreadyInitialized
	}

	c, cancel := context.WithCancel(ctx)
	initialization.cancel = cancel
	initialization.initialized = true

	initialization.iss = iss

	started := make(chan error, 1)
	initialization.started = started
	go initialization.start(c, started)

	initialization.Unlock()

	err := <-started
	if err != nil {
		return err
	}
	if debugEnabled {
		fmt.Printf("kcoidc initialize success: %v\n", iss)
	}

	return nil
}

// Uninitialize uninitializes the global library state.
func Uninitialize() error {
	if debugEnabled {
		fmt.Println("kcoidc uninitialize")
	}

	initialization.Lock()
	defer initialization.Unlock()

	if !initialization.initialized {
		return KCOIDCErrNotInitialized
	}

	initialization.cancel()
	err := <-initialization.started

	initialization.cancel = nil
	initialization.started = nil
	initialization.iss = ""
	initialization.initialized = false
	initialization.ready = nil
	initialization.discovery = nil
	initialization.jwks = nil

	if debugEnabled {
		fmt.Println("kcoidc uninitialize success")
	}

	switch err {
	case context.Canceled:
		return nil
	}
	return err
}

// WaitUntilReady blocks until the initialization is ready or timeout.
func WaitUntilReady(timeout time.Duration) error {
	initialization.RLock()
	if !initialization.initialized {
		initialization.RUnlock()
		return KCOIDCErrNotInitialized
	}
	ready := initialization.ready
	initialization.RUnlock()

	var err error
	if debugEnabled {
		fmt.Println("kcoidc waiting until ready")
		defer func() {
			fmt.Printf("kcoidc finished waiting until ready: %v\n", err)
		}()
	}

	select {
	case <-ready:
	case <-time.After(timeout):
		err = KCOIDCErrTimeout
	}

	return err
}

func (in *initializationData) start(ctx context.Context, started chan error) {
	if debugEnabled {
		fmt.Println("kcoidc start")
		defer func() {
			fmt.Println("kcoidc start has ended")
		}()
	}

	// Close started channel to signal caller that we are started.
	in.Lock()
	if !in.initialized || started != in.started {
		in.Unlock()
		started <- errors.New("start with wrong intialization")
		return
	}

	ready := make(chan struct{})
	in.ready = ready
	in.Unlock()
	started <- nil

	for {
		retry := 60 * time.Second
		if debugEnabled {
			fmt.Println("kcoidc running ...")
		}

		in.RLock()
		if in.initialized && started == in.started {
			iss := in.iss
			in.RUnlock()
			ddoc, err := fetchDiscoveryDocument(ctx, iss)
			if err != nil {
				if debugEnabled {
					fmt.Printf("kcoid discovery failed: %v\n", err)
					retry = 5 * time.Second
				}
			} else {
				jwks, err := fetchJWKSDocument(ctx, ddoc)
				if err != nil {
					if debugEnabled {
						fmt.Printf("kcoid discovery jwks failed: %v\n", err)
						retry = 5 * time.Second
					}
				} else {
					in.Lock()
					if in.initialized && started == in.started {
						in.discovery = ddoc
						in.jwks = jwks
						if debugEnabled {
							fmt.Printf("kcoid ready: %#v, %#v\n", ddoc, jwks)
						}
						close(ready)
					}
					in.Unlock()
				}
			}
		} else {
			in.RUnlock()
		}

		select {
		case <-ctx.Done():
			started <- ctx.Err()
			close(started)
			return
		case <-in.quit:
			close(started)
			return
		case <-time.After(retry):
			// breaks for retry
		}
	}
}
