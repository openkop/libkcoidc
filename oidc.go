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
	"net/http"

	"gopkg.in/square/go-jose.v2"
)

var contentTypeJSONOnly = []string{"application/json"}
var contentTypeJWKSetAndJSON = []string{"application/jwk-set+json", "application/json"}

type oidcDiscoveryDocument struct {
	JWKSUri                   string   `json:"jwks_uri"`
	UserinfoEndpoint          string   `json:"userinfo_endpoint"`
	SupportedSigningAlgValues []string `json:"id_token_signing_alg_values_supported"`
}

func newDiscoveryDocumentFromIssuer(ctx context.Context, client *http.Client, iss string) (*oidcDiscoveryDocument, error) {
	doc := &oidcDiscoveryDocument{}

	return doc, fetchJSON(ctx, client, fmt.Sprintf("%s/.well-known/openid-configuration", iss), nil, contentTypeJSONOnly, doc)
}

type oidcJSONWebKeySet struct {
	*jose.JSONWebKeySet
}

func newoidcJSONWebKeySetFromURL(ctx context.Context, client *http.Client, url string) (*oidcJSONWebKeySet, error) {
	doc := &oidcJSONWebKeySet{
		&jose.JSONWebKeySet{},
	}

	return doc, fetchJSON(ctx, client, url, nil, contentTypeJWKSetAndJSON, doc)
}

type userinfoMap map[string]string
