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
