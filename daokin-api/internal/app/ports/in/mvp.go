package inport

import (
	"context"

	"daokin-api/internal/domain"
)

type AuthCommands interface {
	CreateChallenge(ctx context.Context, wallet string) (domain.AuthChallenge, error)
	VerifyChallenge(ctx context.Context, wallet, nonce, signature string) (domain.AuthSession, error)
	ValidateSession(ctx context.Context, wallet, accessToken string) error
}

type ArtifactCommands interface {
	CreateArtifact(ctx context.Context, creatorWallet, contentHash, contentURI string) (domain.Artifact, error)
}

type ArtifactQueries interface {
	GetArtifact(ctx context.Context, id string) (domain.Artifact, error)
}

type MembershipCommands interface {
	JoinDAO(ctx context.Context, daoID, wallet string) (domain.DAOMembership, error)
	LeaveDAO(ctx context.Context, daoID, wallet, reason string) (domain.DAOMembership, error)
}
