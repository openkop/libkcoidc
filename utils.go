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
	"encoding/json"
)

func popFromMap(m map[string]interface{}, k string) (interface{}, bool) {
	v, ok := m[k]
	if !ok {
		return nil, false
	}

	delete(m, k)
	return v, true
}

func popStringFromMap(m map[string]interface{}, k string) string {
	v, ok := popFromMap(m, k)
	if !ok {
		return ""
	}

	s, _ := v.(string)
	return s
}

func popInt64FromMap(m map[string]interface{}, k string) int64 {
	v, ok := popFromMap(m, k)
	if ok {
		switch vt := v.(type) {
		case float64:
			return int64(vt)
		case json.Number:
			n, _ := vt.Int64()
			return n
		}
	}

	return 0
}
