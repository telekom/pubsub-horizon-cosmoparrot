// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"time"
)

type request struct {
	Time time.Time `json:"time"`
	Path string    `json:"path"`
	// Query is the raw query string without the leading "?", not decoded.
	Query string `json:"query,omitempty"`
	// RequestURI is the request target as received on the request line.
	RequestURI string              `json:"requestUri"`
	Method     string              `json:"method"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Body       json.RawMessage     `json:"body,omitempty"`
	Padding    string              `json:"padding,omitempty"`
}
