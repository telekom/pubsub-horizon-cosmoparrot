// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0
package utils

import (
	"strings"
)

func IsBrowser(userAgent string) bool {
	browserKeywords := []string{"mozilla", "chrome", "safari", "firefox", "edge", "opera"}

	str := strings.ToLower(userAgent)
	for _, keyword := range browserKeywords {
		if strings.Contains(str, keyword) && !strings.Contains(str, "googlebot") {
			return true
		}
	}

	return false
}
