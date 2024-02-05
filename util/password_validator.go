package util

import (
	"regexp"
)

func ValidatePassword(password string) bool {
	if len(password) < 6 || len(password) > 64 {
		return false
	}
	var (
		upperCase = regexp.MustCompile(`[A-Z]`)
		number    = regexp.MustCompile(`[0-9]`)
		special   = regexp.MustCompile(`[@$!%*?&]`)
	)

	if !upperCase.MatchString(password) {
		return false
	}

	if !number.MatchString(password) {
		return false
	}
	if !special.MatchString(password) {
		return false
	}
	return true
}
