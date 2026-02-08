package domain

import (
	"strings"
	"time"
)

type Artifact struct {
	ID            string    `json:"id"`
	CreatorWallet string    `json:"creator_wallet"`
	ContentHash   string    `json:"content_hash"`
	ContentURI    string    `json:"content_uri"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewArtifact(creatorWallet, contentHash, contentURI string, now time.Time) (Artifact, error) {
	creatorWallet = strings.TrimSpace(creatorWallet)
	contentHash = strings.TrimSpace(contentHash)
	contentURI = strings.TrimSpace(contentURI)

	if creatorWallet == "" || contentHash == "" || contentURI == "" {
		return Artifact{}, ErrInvalidArgument
	}

	return Artifact{
		CreatorWallet: creatorWallet,
		ContentHash:   contentHash,
		ContentURI:    contentURI,
		CreatedAt:     now.UTC(),
	}, nil
}
