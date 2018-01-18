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

// #define KOIDC_API
import "C"

import (
	"fmt"
)

// KCOIDCErr is the error type as used by this library.
type KCOIDCErr uint64

func (err KCOIDCErr) Error() string {
	return fmt.Sprintf("%s (KCOIDC:0x%x)", KCOIDCErrText(err), uint64(err))
}

// Errors as defined by this library.
const (
	KCOIDCErrNone              = iota
	KCOIDCErrUnknown KCOIDCErr = (1 << 8) | iota
	KCOIDCErrAlreadyInitialized
	KCOIDCErrNotInitialized
	KCOIDCErrTimeout
	KCOIDCErrTokenUnexpectedSigningMethod
	KCOIDCErrTokenMalformed
	KCOIDCErrTokenExpiredOrNotValidYet
	KCOIDCErrTokenUnknownKey
	KCOIDCErrTokenInvalidSignature
	KCOIDCErrTokenValidationFailed
)

// KCOIDCSuccess is the success response as returned by this library.
const KCOIDCSuccess = KCOIDCErrNone

// KCOIDCErrTextMap maps erros to readable names.
var KCOIDCErrTextMap = map[KCOIDCErr]string{
	KCOIDCErrUnknown:                      "Unknown",
	KCOIDCErrAlreadyInitialized:           "Already Initialized",
	KCOIDCErrNotInitialized:               "Not Initialized",
	KCOIDCErrTimeout:                      "Timeout",
	KCOIDCErrTokenUnexpectedSigningMethod: "Unexpected Token Signing Method",
	KCOIDCErrTokenMalformed:               "Malformed Token",
	KCOIDCErrTokenExpiredOrNotValidYet:    "Token Expired Or Not Valid Yet",
	KCOIDCErrTokenUnknownKey:              "Unknown Token Key",
	KCOIDCErrTokenInvalidSignature:        "Invalid Token Signature",
	KCOIDCErrTokenValidationFailed:        "Token Validation Failed",
}

// KCOIDCErrText returns a text for the KCOIDCErr. It returns the empty string
// if the code is unknown.
func KCOIDCErrText(code KCOIDCErr) string {
	text := KCOIDCErrTextMap[code]
	return text
}

func returnKCOIDCErrorOrUnknown(err error) C.ulonglong {
	switch e := err.(type) {
	case KCOIDCErr:
		return C.ulonglong(e)
	default:
		if debugEnabled {
			fmt.Printf("kcoid unknown error: %s\n", err)
		}
		return C.ulonglong(KCOIDCErrUnknown)
	}
}
