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
	"github.com/dgrijalva/jwt-go"
)

// Token claims used by Kopano Konnect.
const (
	IsAccessTokenClaim  = "kc.isAccessToken"
	IsRefreshTokenClaim = "kc.isRefreshToken"

	IdentityClaim         = "kc.identity"
	IdentifiedUserIDClaim = "kc.i.id"
)

// Token types as int.
const (
	TokenTypeStandard  int = 0
	TokenTypeKCAccess  int = 1
	TokenTypeKCRefresh int = 2
)

// ExtraClaimsWithType is a MapClaims with a specific type.
type ExtraClaimsWithType jwt.MapClaims

// Valid satisfies the jwt.Claims interface.
func (claims *ExtraClaimsWithType) Valid() error {
	return nil
}

// KCTokenType returns the numeric type of the accociated claims.
func (claims *ExtraClaimsWithType) KCTokenType() int {
	if v, _ := (*claims)[IsAccessTokenClaim].(bool); v {
		return TokenTypeKCAccess
	}
	if v, _ := (*claims)[IsRefreshTokenClaim].(bool); v {
		return TokenTypeKCRefresh
	}

	return TokenTypeStandard
}

// SplitStandardClaimsFromMapClaims removes all JWT standard claims from the
// provided map claims and returns them.
func SplitStandardClaimsFromMapClaims(claims *ExtraClaimsWithType) (*jwt.StandardClaims, error) {
	std := &jwt.StandardClaims{
		Audience:  popStringFromMap(*claims, "aud"),
		ExpiresAt: popInt64FromMap(*claims, "exp"),
		Id:        popStringFromMap(*claims, "jti"),
		IssuedAt:  popInt64FromMap(*claims, "iat"),
		Issuer:    popStringFromMap(*claims, "iss"),
		NotBefore: popInt64FromMap(*claims, "nbf"),
		Subject:   popStringFromMap(*claims, "sub"),
	}

	return std, nil
}

// AuthenticatedUserIDFromClaims extracts extra Kopano Connect identified claims
// from the provided extra claims.
func AuthenticatedUserIDFromClaims(claims *ExtraClaimsWithType) (string, bool) {
	if identityClaims, _ := (*claims)[IdentityClaim].(map[string]interface{}); identityClaims != nil {
		authenticatedUserID, _ := identityClaims[IdentifiedUserIDClaim].(string)
		return authenticatedUserID, true
	}

	return "", false
}
