package domain

import "time"

type Artifact struct {
	ID            string    `json:"id"`
	CreatorWallet string    `json:"creator_wallet"`
	ContentHash   string    `json:"content_hash"`
	ContentURI    string    `json:"content_uri"`
	CreatedAt     time.Time `json:"created_at"`
}
