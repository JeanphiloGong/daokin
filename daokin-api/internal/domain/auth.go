package domain

import (
	"strings"
	"time"
)

type AuthChallenge struct {
	Wallet    string    `json:"wallet"`
	Nonce     string    `json:"nonce"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AuthSession struct {
	Wallet      string    `json:"wallet"`
	AccessToken string    `json:"access_token"`
	VerifiedAt  time.Time `json:"verified_at"`
}

func NewAuthChallenge(wallet, nonce, message string, now time.Time, ttl time.Duration) (AuthChallenge, error) {
	wallet = strings.TrimSpace(wallet)
	nonce = strings.TrimSpace(nonce)
	message = strings.TrimSpace(message)
	if wallet == "" || nonce == "" || message == "" || ttl <= 0 {
		return AuthChallenge{}, ErrInvalidArgument
	}

	now = now.UTC()
	return AuthChallenge{
		Wallet:    wallet,
		Nonce:     nonce,
		Message:   message,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}, nil
}

func (c AuthChallenge) Verify(nonce string, now time.Time) error {
	if strings.TrimSpace(nonce) == "" {
		return ErrInvalidArgument
	}
	if c.Nonce != strings.TrimSpace(nonce) {
		return ErrChallengeInvalid
	}
	if now.UTC().After(c.ExpiresAt) {
		return ErrChallengeInvalid
	}
	return nil
}

func NewAuthSession(wallet, token string, now time.Time) (AuthSession, error) {
	wallet = strings.TrimSpace(wallet)
	token = strings.TrimSpace(token)
	if wallet == "" || token == "" {
		return AuthSession{}, ErrInvalidArgument
	}

	return AuthSession{
		Wallet:      wallet,
		AccessToken: token,
		VerifiedAt:  now.UTC(),
	}, nil
}
