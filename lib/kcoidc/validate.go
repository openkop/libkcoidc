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
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// ValidateTokenString validates the provided token string value and returns
// the subject as found in the claims.
func ValidateTokenString(tokenString string) (string, error) {
	initialization.RLock()
	ddoc := initialization.discovery
	jwks := initialization.jwks
	initialization.RUnlock()
	if ddoc == nil || jwks == nil {
		return "", KCOIDCErrNotInitialized
	}

	if debugEnabled {
		fmt.Printf("kcoidc validate token string: %s\n", tokenString)
	}

	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if debugEnabled {
			fmt.Printf("kcoidc validate token header: %#v\n", token.Header)
		}

		supportedAlg := false
		for _, alg := range ddoc.SupportedSigningAlgValues {
			if token.Method.Alg() == alg {
				supportedAlg = true
				break
			}
		}
		if !supportedAlg {
			return nil, KCOIDCErrTokenUnexpectedSigningMethod
		}

		kid, _ := (token.Header["kid"].(string))
		keys := jwks.Key(kid)
		if keys == nil || len(keys) == 0 {
			return nil, KCOIDCErrTokenUnknownKey
		}

		key := keys[0]
		if debugEnabled {
			fmt.Printf("kcoidc validate token key: %#v\n", key.Key)
		}

		return key.Key, nil
	})

	if debugEnabled {
		fmt.Printf("kcoidc validate token result: %v, %#v, %s\n", claims.Subject, token.Valid, err)
	}

	if token != nil && token.Valid {
		return claims.Subject, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			err = KCOIDCErrTokenMalformed
		} else if ve.Errors&(jwt.ValidationErrorSignatureInvalid|jwt.ValidationErrorUnverifiable) != 0 {
			err = KCOIDCErrTokenInvalidSignature
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			err = KCOIDCErrTokenExpiredOrNotValidYet
		} else {
			err = KCOIDCErrTokenValidationFailed
		}
	}

	if debugEnabled {
		fmt.Printf("kcoid validate token resulted in validation failure: %s\n", err)
	}
	return claims.Subject, err
}
