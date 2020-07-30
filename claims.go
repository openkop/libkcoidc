/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
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
	IdentifiedUserIsGuest = "kc.i.guest"

	AuthorizedScopesClaim = "kc.authorizedScopes"
	AuthorizedClaimsClaim = "kc.authorizedClaims"
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
// from the provided extra claims, returning the authenticated user id.
func AuthenticatedUserIDFromClaims(claims *ExtraClaimsWithType) (string, bool) {
	if identityClaims, _ := (*claims)[IdentityClaim].(map[string]interface{}); identityClaims != nil {
		if authenticatedUserID, _ := identityClaims[IdentifiedUserIDClaim].(string); authenticatedUserID != "" {
			return authenticatedUserID, true
		}
	}

	return "", false
}

// AuthenticatedUserIsGuest extract extra Kopano Connect identified claims from
// the provided extra claims, returning if the claims are for a guest or not.
func AuthenticatedUserIsGuest(claims *ExtraClaimsWithType) bool {
	if identityClaims, _ := (*claims)[IdentityClaim].(map[string]interface{}); identityClaims != nil {
		isGuest, _ := identityClaims[IdentifiedUserIsGuest].(bool)
		return isGuest
	}

	return false
}

// AuthorizedScopesFromClaims returns the authorized scopes as bool map from
// the provided extra claims.
func AuthorizedScopesFromClaims(claims *ExtraClaimsWithType) map[string]bool {
	if authorizedScopes, _ := (*claims)[AuthorizedScopesClaim].([]interface{}); authorizedScopes != nil {
		authorizedScopesMap := make(map[string]bool)
		for _, scope := range authorizedScopes {
			authorizedScopesMap[scope.(string)] = true
		}

		return authorizedScopesMap
	}

	return nil
}

// AuthorizedClaimsFromClaims returns the authorized claims as map from the
// provided extra claims.
func AuthorizedClaimsFromClaims(claims *ExtraClaimsWithType) map[string]interface{} {
	authorizedClaims, _ := (*claims)[AuthorizedClaimsClaim].(map[string]interface{})

	return authorizedClaims
}

// RequireScopesInClaims returns nil if all the provided scopes are found in
// the provided claims. Otherwise an error is returned.
func RequireScopesInClaims(claims *ExtraClaimsWithType, requiredScopes []string) error {
	if len(requiredScopes) == 0 {
		return nil
	}

	authorizedScopes := AuthorizedScopesFromClaims(claims)
	missingScopes := make([]string, 0)
	for _, scope := range requiredScopes {
		if ok, _ := authorizedScopes[scope]; !ok {
			missingScopes = append(missingScopes, scope)
		}
	}
	if len(missingScopes) == 0 {
		return nil
	}

	return ErrStatusMissingRequiredScope
}
