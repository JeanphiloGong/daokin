package domain

import (
	"crypto/sha256"
	"time"
)

type Session struct {
	TokenHash [32]byte
	UserID    string
	ExpiresAt time.Time
}

func TokenHash(token string) [32]byte { return sha256.Sum256([]byte(token)) }
func NewSession(userID, token string, now time.Time) Session {
	return Session{TokenHash: TokenHash(token), UserID: userID, ExpiresAt: now.UTC().Truncate(time.Microsecond).Add(24 * time.Hour)}
}
func (s Session) Active(now time.Time) bool { return now.Before(s.ExpiresAt) }
