// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

type response struct {
	Path    string              `json:"path"`
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    any                 `json:"body,omitempty"`
}
