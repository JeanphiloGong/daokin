package outport

import (
	"context"
	"time"

	"daokin-api/internal/domain"
)

type AuthRepository interface {
	SaveChallenge(ctx context.Context, challenge domain.AuthChallenge) error
	GetChallenge(ctx context.Context, wallet string) (domain.AuthChallenge, error)
	DeleteChallenge(ctx context.Context, wallet string) error
	SaveSession(ctx context.Context, session domain.AuthSession) error
}

type ArtifactRepository interface {
	Create(ctx context.Context, artifact domain.Artifact) (domain.Artifact, error)
	GetByID(ctx context.Context, id string) (domain.Artifact, error)
}

type MembershipRepository interface {
	Join(ctx context.Context, membership domain.DAOMembership) (domain.DAOMembership, error)
	Leave(ctx context.Context, daoID, wallet string, leftAt time.Time, reason string) (domain.DAOMembership, error)
}
