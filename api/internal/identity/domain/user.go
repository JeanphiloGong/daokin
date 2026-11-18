// defines a user with an ID and wallet address, plus two functions for validating wallet input and constructing a user.
package domain

import (
	"net/mail"
	"strings"
	"unicode"

	"github.com/JeanphiloGong/daokin/api/internal/problem"
)

type User struct {
	ID    string
	Email string
}

func NormalizeEmail(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 254 || strings.ContainsFunc(value, func(r rune) bool { return r > unicode.MaxASCII }) {
		return "", problem.Invalid
	}
	parsed, err := mail.ParseAddress(value)
	if err != nil || parsed.Address != value {
		return "", problem.Invalid
	}
	return strings.ToLower(value), nil
}
