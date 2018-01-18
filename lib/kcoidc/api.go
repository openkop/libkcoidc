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

import (
	"C"
)

//export kcoidc_validate_token_s
func kcoidc_validate_token_s(tokenCString *C.char) (*C.char, bool) {
	subject, err := ValidateTokenString(C.GoString(tokenCString))
	if err != nil {
		return C.CString(subject), false
	}
	return C.CString(subject), true
}
