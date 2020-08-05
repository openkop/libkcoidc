/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kcoidc

import (
	"fmt"
)

// ErrStatus is the Error type as used by kcoidc.
type ErrStatus uint64

func (errStatus ErrStatus) Error() string {
	return fmt.Sprintf("%s (:0x%x)", ErrStatusText(errStatus), uint64(errStatus))
}

// ErrStatusors as defined by this library.
const (
	ErrStatusNone              = iota
	ErrStatusUnknown ErrStatus = iota + (1 << 8)
	ErrStatusInvalidIss
	ErrStatusAlreadyInitialized
	ErrStatusNotInitialized
	ErrStatusTimeout
	ErrStatusTokenUnexpectedSigningMethod
	ErrStatusTokenMalformed
	ErrStatusTokenExpiredOrNotValidYet
	ErrStatusTokenUnknownKey
	ErrStatusTokenInvalidSignature
	ErrStatusTokenValidationFailed
	ErrStatusClosed
	ErrStatusWrongInitialization
	ErrStatusMissingRequiredScope
)

// StatusSuccess is the success response as returned by this library.
const StatusSuccess = ErrStatusNone

// ErrStatusTextMap maps ErrStatusos to readable names.
var ErrStatusTextMap = map[ErrStatus]string{
	ErrStatusUnknown:                      "Unknown",
	ErrStatusInvalidIss:                   "Invalid Issuer Identifier Value",
	ErrStatusAlreadyInitialized:           "Already Initialized",
	ErrStatusNotInitialized:               "Not Initialized",
	ErrStatusTimeout:                      "Timeout",
	ErrStatusTokenUnexpectedSigningMethod: "Unexpected Token Signing Method",
	ErrStatusTokenMalformed:               "Malformed Token",
	ErrStatusTokenExpiredOrNotValidYet:    "Token Expired Or Not Valid Yet",
	ErrStatusTokenUnknownKey:              "Unknown Token Key",
	ErrStatusTokenInvalidSignature:        "Invalid Token Signature",
	ErrStatusTokenValidationFailed:        "Token Validation Failed",
	ErrStatusClosed:                       "Is Closed",
	ErrStatusWrongInitialization:          "Wrong Initialization",
	ErrStatusMissingRequiredScope:         "Missing required scope",
}

// ErrStatusText returns a text for the ErrStatus. It returns the empty string
// if the code is unknown.
func ErrStatusText(code ErrStatus) string {
	text := ErrStatusTextMap[code]
	return text
}
