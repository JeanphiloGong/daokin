package attribution

import (
	"context"
	"strings"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type QueryService struct {
	permissions outport.PermissionRepository
	transfers   outport.TransferRepository
}

func NewQueryService(
	permissions outport.PermissionRepository,
	transfers outport.TransferRepository,
) *QueryService {
	return &QueryService{
		permissions: permissions,
		transfers:   transfers,
	}
}

func (s *QueryService) GetArtifactTrail(ctx context.Context, artifactID string) (domain.AttributionTrail, error) {
	artifactID = strings.TrimSpace(artifactID)
	if artifactID == "" {
		return domain.AttributionTrail{}, domain.ErrInvalidArgument
	}

	permissions, err := s.permissions.ListPermissionsByArtifact(ctx, artifactID)
	if err != nil {
		return domain.AttributionTrail{}, err
	}
	transfers, err := s.transfers.ListTransfersByArtifact(ctx, artifactID)
	if err != nil {
		return domain.AttributionTrail{}, err
	}

	return domain.AttributionTrail{
		ArtifactID:   artifactID,
		Permissions:  permissions,
		Transfers:    transfers,
		TotalRecords: len(permissions) + len(transfers),
	}, nil
}
