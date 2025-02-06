package utils

import (
	"strings"
)

func IsBrowser(userAgent string) bool {
	browserKeywords := []string{"mozilla", "chrome", "safari", "firefox", "edge", "opera"}

	str := strings.ToLower(userAgent)
	for _, keyword := range browserKeywords {
		if strings.Contains(str, keyword) {
			return true
		}
	}

	return false
}
