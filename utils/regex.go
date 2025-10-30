package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var reUsername = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]{3,29}$`)
var rePhone = regexp.MustCompile(`^[0-9]{10}$`)
var reEmail = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

func IsValidUsername(s string) bool {
	return reUsername.MatchString(s)
}

func IsValidPhone(s string) bool {
	return rePhone.MatchString(s)
}

func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func IsValidEmail(s string) bool {
	return reEmail.MatchString(s)
}

func ValidatePassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case strings.ContainsRune("@$!%*?&", c):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}