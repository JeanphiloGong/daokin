package domain

import "time"

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
