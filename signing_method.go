/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kcoidc

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
)

func init() {
	fixupRSAPSSSaltLength()
}

func fixupRSAPSSSaltLength() {
	for _, name := range []string{"PS256", "PS386", "PS512"} {
		signingMethod := jwt.GetSigningMethod(name)
		if signingMethodRSAPSS, ok := signingMethod.(*jwt.SigningMethodRSAPSS); ok {
			// NOTE(longsleep): Ensure to use same salt length the hash size.
			// See https://www.ietf.org/mail-archive/web/jose/current/msg02901.html for
			// reference and https://github.com/dgrijalva/jwt-go/issues/285 for
			// the issue in upstream jwt-go.
			signingMethodRSAPSS.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
		}
	}
}
