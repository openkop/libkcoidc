/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
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
