// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

type response struct {
	Time    time.Time           `json:"time"`
	Path    string              `json:"path"`
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    any                 `json:"body,omitempty"`
}
