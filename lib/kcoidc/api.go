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
	"context"
	"time"
)

//export kcoidc_initialize
func kcoidc_initialize(issCString *C.char) C.ulonglong {
	err := Initialize(context.Background(), C.GoString(issCString))
	if err != nil {
		return returnKCOIDCErrorOrUnknown(err)
	}
	return KCOIDCSuccess
}

//export kcoidc_wait_untill_ready
func kcoidc_wait_untill_ready(timeout C.ulonglong) C.ulonglong {
	err := WaitUntilReady(time.Duration(timeout) * time.Second)
	if err != nil {
		return returnKCOIDCErrorOrUnknown(err)
	}
	return KCOIDCSuccess
}

//export kcoidc_insecure_skip_verify
func kcoidc_insecure_skip_verify(enableInsecure C.int) C.ulonglong {
	err := InsecureSkipVerify(enableInsecure == 1)
	if err != nil {
		return returnKCOIDCErrorOrUnknown(err)
	}
	return KCOIDCSuccess
}

//export kcoidc_validate_token_s
func kcoidc_validate_token_s(tokenCString *C.char) (*C.char, C.ulonglong) {
	subject, err := ValidateTokenString(C.GoString(tokenCString))
	if err != nil {
		return C.CString(subject), returnKCOIDCErrorOrUnknown(err)
	}
	return C.CString(subject), KCOIDCSuccess
}

//export kcoidc_uninitialize
func kcoidc_uninitialize() C.ulonglong {
	err := Uninitialize()
	if err != nil {
		return returnKCOIDCErrorOrUnknown(err)
	}
	return KCOIDCSuccess
}
