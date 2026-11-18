package domain

import (
	"testing"
	"time"
)

func TestEmailIdentity(t *testing.T) {
	email, err := NormalizeEmail("  Lin+Writing@Example.com  ")
	if err != nil || email != "lin+writing@example.com" {
		t.Fatalf("email = %q, error = %v", email, err)
	}
	for _, bad := range []string{"lin", "Lin <lin@example.com>", "林@example.com"} {
		if _, err := NormalizeEmail(bad); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}
func TestSessionExpiry(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := NewSession("user_a", "synthetic-token", now)
	if !s.Active(now.Add(23*time.Hour)) || s.Active(now.Add(24*time.Hour)) || s.UserID != "user_a" || s.TokenHash != TokenHash("synthetic-token") {
		t.Fatal("session lost the subject, digest, or fixed deadline")
	}
}
