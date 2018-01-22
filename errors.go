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
	"fmt"
)

// ErrStatus is the Error type as used by kcoidc.
type ErrStatus uint64

func (ErrStatus ErrStatus) Error() string {
	return fmt.Sprintf("%s (:0x%x)", ErrStatusText(ErrStatus), uint64(ErrStatus))
}

// ErrStatusors as defined by this library.
const (
	ErrStatusNone              = iota
	ErrStatusUnknown ErrStatus = (1 << 8) | iota
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
}

// ErrStatusText returns a text for the ErrStatus. It returns the empty string
// if the code is unknown.
func ErrStatusText(code ErrStatus) string {
	text := ErrStatusTextMap[code]
	return text
}
