/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
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
