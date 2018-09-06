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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func fetchJSON(ctx context.Context, client *http.Client, url string, headers http.Header, validContentTypes []string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if client == nil {
		client = http.DefaultClient
	}
	if headers != nil {
		req.Header = headers
	}

	req = req.WithContext(ctx)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if len(validContentTypes) > 0 {
		contentType := strings.SplitN(response.Header.Get("Content-Type"), ";", 2)[0]
		valid := false
		for _, ct := range validContentTypes {
			if ct == contentType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unexpected response content-type: %s", contentType)
		}
	}

	return json.NewDecoder(response.Body).Decode(target)
}
