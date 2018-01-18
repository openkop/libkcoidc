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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type oidcDiscoveryDocument struct {
	JWKSUri                   string   `json:"jwks_uri"`
	SupportedSigningAlgValues []string `json:"id_token_signing_alg_values_supported"`
}

func fetchDiscoveryDocument(ctx context.Context, iss string) (*oidcDiscoveryDocument, error) {
	url := fmt.Sprintf("%s/.well-known/openid-configuration", iss)
	ddoc := &oidcDiscoveryDocument{}

	err := fetchJSON(ctx, url, ddoc)
	if err != nil {
		return nil, err
	}

	return ddoc, nil
}

func fetchJWKSDocument(ctx context.Context, ddoc *oidcDiscoveryDocument) (*jose.JSONWebKeySet, error) {
	url := ddoc.JWKSUri
	jwks := &jose.JSONWebKeySet{}

	err := fetchJSON(ctx, url, jwks)
	if err != nil {
		return nil, err
	}

	return jwks, nil
}

func fetchJSON(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req = req.WithContext(c)

	response, err := initialization.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}
