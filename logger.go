/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kcoidc

// A Logger defines a simple logging interface for pluggable loggers used by
// this module.
type Logger interface {
	Printf(string, ...interface{})
}

// DefaultLogger is the logger used by this library if no other is explicitly
// specified.
var DefaultLogger Logger = nil
