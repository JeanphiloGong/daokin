package outport

import (
	"context"
	"time"

	"daokin-api/internal/domain"
)

type AuthRepository interface {
	SaveChallenge(ctx context.Context, challenge domain.AuthChallenge) error
	GetChallenge(ctx context.Context, wallet string) (domain.AuthChallenge, error)
	ConsumeChallenge(ctx context.Context, wallet string) (domain.AuthChallenge, error)
	DeleteChallenge(ctx context.Context, wallet string) error
	SaveSession(ctx context.Context, session domain.AuthSession) error
	GetSession(ctx context.Context, wallet string) (domain.AuthSession, error)
}

type SignatureVerifier interface {
	VerifyWalletSignature(ctx context.Context, wallet, message, signature string) error
}

type ArtifactRepository interface {
	Create(ctx context.Context, artifact domain.Artifact) (domain.Artifact, error)
	GetByID(ctx context.Context, id string) (domain.Artifact, error)
	ListByCreator(ctx context.Context, wallet string) ([]domain.Artifact, error)
}

type MembershipRepository interface {
	Join(ctx context.Context, membership domain.DAOMembership) (domain.DAOMembership, error)
	Leave(ctx context.Context, daoID, wallet string, leftAt time.Time, reason string) (domain.DAOMembership, error)
	ListByWalletMemberships(ctx context.Context, wallet string) ([]domain.DAOMembership, error)
}

type PermissionRepository interface {
	Grant(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	Revoke(ctx context.Context, permissionID, granterWallet string, revokedAt time.Time, reason string) (domain.Permission, error)
	HasActive(ctx context.Context, artifactID, granteeWallet, scope string) (bool, error)
	ListPermissionsByArtifact(ctx context.Context, artifactID string) ([]domain.Permission, error)
	ListByWalletPermissions(ctx context.Context, wallet string) ([]domain.Permission, error)
}

type TransferRepository interface {
	CreateTransfer(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error)
	ListTransfersByArtifact(ctx context.Context, artifactID string) ([]domain.Transfer, error)
	ListByWalletTransfers(ctx context.Context, wallet string) ([]domain.Transfer, error)
}
