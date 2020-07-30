/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kcoidc

import "testing"

func TestErrors(t *testing.T) {
	for err := range ErrStatusTextMap {
		t.Logf("%d: %s", err, err)
	}
}
