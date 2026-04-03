package permission

import (
	"context"
	"strings"
	"time"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type CommandService struct {
	permissions outport.PermissionRepository
	artifacts   outport.ArtifactRepository
	now         func() time.Time
}

func NewCommandService(
	permissions outport.PermissionRepository,
	artifacts outport.ArtifactRepository,
) *CommandService {
	return &CommandService{
		permissions: permissions,
		artifacts:   artifacts,
		now:         time.Now,
	}
}

func (s *CommandService) GrantPermission(
	ctx context.Context,
	artifactID string,
	granterWallet string,
	granteeWallet string,
	scope string,
) (domain.Permission, error) {
	artifactID = strings.TrimSpace(artifactID)
	granterWallet = strings.TrimSpace(granterWallet)

	artifact, err := s.artifacts.GetByID(ctx, artifactID)
	if err != nil {
		return domain.Permission{}, err
	}
	if !strings.EqualFold(artifact.CreatorWallet, granterWallet) {
		return domain.Permission{}, domain.ErrForbidden
	}

	permission, err := domain.NewPermission(artifactID, granterWallet, granteeWallet, scope, s.now())
	if err != nil {
		return domain.Permission{}, err
	}
	return s.permissions.Grant(ctx, permission)
}

func (s *CommandService) RevokePermission(
	ctx context.Context,
	permissionID string,
	granterWallet string,
	reason string,
) (domain.Permission, error) {
	permissionID = strings.TrimSpace(permissionID)
	granterWallet = strings.TrimSpace(granterWallet)
	if permissionID == "" || granterWallet == "" {
		return domain.Permission{}, domain.ErrInvalidArgument
	}
	return s.permissions.Revoke(ctx, permissionID, granterWallet, s.now(), reason)
}
