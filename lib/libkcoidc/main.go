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
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// ValidateToken validates a raw JWT token string.
//export ValidateToken
func ValidateToken(tokenRawString *C.char) (*C.char, bool) {
	tokenString := C.GoString(tokenRawString)

	fmt.Printf("tokenString: %s\n", tokenString)

	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	fmt.Printf("token validation result: %#v, %#v, %v\n", token, claims, err)

	return C.CString(claims.Subject), token.Valid
}

func main() {}
