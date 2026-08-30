package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"unicode/utf8"
)

const (
	MinPasswordLength = 6
	MaxPasswordLength = 64
)

var (
	reUpper = regexp.MustCompile(`[A-Z]`)
	reLower = regexp.MustCompile(`[a-z]`)
	reDigit = regexp.MustCompile(`[0-9]`)
)

// ValidatePassword enforces the password policy:
//   - length between MinPasswordLength and MaxPasswordLength
//   - must contain at least one uppercase letter
//   - must contain at least one lowercase letter
//   - must contain at least one digit
//
// It returns false and a human-readable reason when the password is rejected.
func ValidatePassword(pwd string) (bool, string) {
	if utf8.RuneCountInString(pwd) < MinPasswordLength {
		return false, "密码长度至少 6 位"
	}
	if utf8.RuneCountInString(pwd) > MaxPasswordLength {
		return false, "密码长度不能超过 64 位"
	}
	if !reUpper.MatchString(pwd) {
		return false, "密码必须包含至少一个大写字母"
	}
	if !reLower.MatchString(pwd) {
		return false, "密码必须包含至少一个小写字母"
	}
	if !reDigit.MatchString(pwd) {
		return false, "密码必须包含至少一个数字"
	}
	return true, ""
}

// UsernameFromPassword derives a stable username from the password itself, so
// the same password always identifies the same account.
func UsernameFromPassword(pwd string) string {
	sum := sha256.Sum256([]byte(pwd))
	return "u" + hex.EncodeToString(sum[:])[:32]
}
