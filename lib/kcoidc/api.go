/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

/*
#define KCOIDC_API 1
#define KCOIDC_API_MINOR 1

#define KCOIDC_VERSION (KCOIDC_API * 10000 + KCOIDC_API_MINOR * 100)

// Token types as defined by kcoidc in claims.go, made usable from C.
static int const KCOIDC_TOKEN_TYPE_STANDARD = 0;
static int const KCOIDC_TOKEN_TYPE_KCACCESS = 1;
static int const KCOIDC_TOKEN_TYPE_KCRERESH = 2;
*/
import "C"

import (
	"context"
	"encoding/json"
	"time"

	"stash.kopano.io/kc/libkcoidc"
)

//export kcoidc_version
func kcoidc_version() *C.char {
	return C.CString(Version())
}

//export kcoidc_build_date
func kcoidc_build_date() *C.char {
	return C.CString(BuildDate())
}

//export kcoidc_initialize
func kcoidc_initialize(issCString *C.char) C.ulonglong {
	err := Initialize(context.Background(), C.GoString(issCString))
	if err != nil {
		return asKnownErrorOrUnknown(err)
	}
	return kcoidc.StatusSuccess
}

//export kcoidc_wait_until_ready
func kcoidc_wait_until_ready(timeout C.ulonglong) C.ulonglong {
	err := WaitUntilReady(time.Duration(timeout) * time.Second)
	if err != nil {
		return asKnownErrorOrUnknown(err)
	}
	return kcoidc.StatusSuccess
}

//export kcoidc_insecure_skip_verify
func kcoidc_insecure_skip_verify(enableInsecure C.int) C.ulonglong {
	err := InsecureSkipVerify(enableInsecure == 1)
	if err != nil {
		return asKnownErrorOrUnknown(err)
	}
	return kcoidc.StatusSuccess
}

//export kcoidc_validate_token_s
func kcoidc_validate_token_s(tokenCString *C.char) (*C.char, C.ulonglong, C.int, *C.char, *C.char) {
	var standardClaimsBytes []byte
	var extraClaimsBytes []byte
	tokenType := kcoidc.TokenTypeStandard
	subject, standardClaims, extraClaims, err := ValidateTokenString(C.GoString(tokenCString))
	if standardClaims != nil {
		// Encode to JSON
		standardClaimsBytes, _ = json.Marshal(standardClaims)
	}
	if extraClaims != nil {
		// Encode to JSON
		extraClaimsBytes, _ = json.Marshal(extraClaims)
		tokenType = extraClaims.KCTokenType()
	}
	if err != nil {
		return C.CString(subject), asKnownErrorOrUnknown(err), C.int(tokenType), C.CString(string(standardClaimsBytes)), C.CString(string(extraClaimsBytes))
	}
	return C.CString(subject), kcoidc.StatusSuccess, C.int(tokenType), C.CString(string(standardClaimsBytes)), C.CString(string(extraClaimsBytes))
}

//export kcoidc_validate_token_and_require_scope_s
func kcoidc_validate_token_and_require_scope_s(tokenCString *C.char, requiredScopeCString *C.char) (*C.char, C.ulonglong, C.int, *C.char, *C.char) {
	var standardClaimsBytes []byte
	var extraClaimsBytes []byte
	tokenType := kcoidc.TokenTypeStandard
	subject, standardClaims, extraClaims, err := ValidateTokenStringAndRequireClaim(C.GoString(tokenCString), C.GoString(requiredScopeCString))
	if standardClaims != nil {
		// Encode to JSON
		standardClaimsBytes, _ = json.Marshal(standardClaims)
	}
	if extraClaims != nil {
		// Encode to JSON
		extraClaimsBytes, _ = json.Marshal(extraClaims)
		tokenType = extraClaims.KCTokenType()
	}
	if err != nil {
		return C.CString(subject), asKnownErrorOrUnknown(err), C.int(tokenType), C.CString(string(standardClaimsBytes)), C.CString(string(extraClaimsBytes))
	}
	return C.CString(subject), kcoidc.StatusSuccess, C.int(tokenType), C.CString(string(standardClaimsBytes)), C.CString(string(extraClaimsBytes))
}

//export kcoidc_fetch_userinfo_with_accesstoken_s
func kcoidc_fetch_userinfo_with_accesstoken_s(tokenCString *C.char) (*C.char, C.ulonglong) {
	userinfo, err := FetchUserinfoWithAccesstokenString(C.GoString(tokenCString))
	if err != nil {
		return nil, asKnownErrorOrUnknown(err)
	}

	// Encode to JSON
	res, err := json.Marshal(userinfo)
	if err != nil {
		return nil, asKnownErrorOrUnknown(err)
	}

	return C.CString(string(res)), kcoidc.StatusSuccess
}

//export kcoidc_uninitialize
func kcoidc_uninitialize() C.ulonglong {
	err := Uninitialize()
	if err != nil {
		return asKnownErrorOrUnknown(err)
	}
	return kcoidc.StatusSuccess
}
