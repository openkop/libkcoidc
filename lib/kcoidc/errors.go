/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

import (
	"C"
	"fmt"

	"stash.kopano.io/kc/libkcoidc"
)

func asKnownErrorOrUnknown(err error) C.ulonglong {
	switch e := err.(type) {
	case kcoidc.ErrStatus:
		return C.ulonglong(e)
	default:
		if debug {
			fmt.Printf("kcoidc-c unknown error: %s\n", err)
		}
		return C.ulonglong(kcoidc.ErrStatusUnknown)
	}
}
