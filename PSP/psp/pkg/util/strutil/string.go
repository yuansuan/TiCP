package strutil

import (
	"strings"
)

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

func DefaultStr(s, defaultStr string) string {
	if IsEmpty(s) {
		return defaultStr
	}

	return s
}
