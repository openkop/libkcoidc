/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kcoidc

type logger interface {
	Printf(string, ...interface{})
}

// DefaultLogger is the logger used by this library if no other is explicitly
// specified.
var DefaultLogger logger = nil
