package artifact

import (
	"context"
	"strings"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type QueryService struct {
	repo outport.ArtifactRepository
}

func NewQueryService(repo outport.ArtifactRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) GetArtifact(ctx context.Context, id string) (domain.Artifact, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.Artifact{}, domain.ErrInvalidArgument
	}
	return s.repo.GetByID(ctx, id)
}
