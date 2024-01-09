package utils

import "strings"

// Returns s with the first ascii letter mapped to their lower case.
func ToLowerFirst(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

// Returns s with the first ascii letter mapped to their upper case.
func ToUpperFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
