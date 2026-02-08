package artifact

import (
	"context"
	"time"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type CommandService struct {
	repo outport.ArtifactRepository
	now  func() time.Time
}

func NewCommandService(repo outport.ArtifactRepository) *CommandService {
	return &CommandService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *CommandService) CreateArtifact(
	ctx context.Context,
	creatorWallet string,
	contentHash string,
	contentURI string,
) (domain.Artifact, error) {
	artifact, err := domain.NewArtifact(creatorWallet, contentHash, contentURI, s.now())
	if err != nil {
		return domain.Artifact{}, err
	}
	return s.repo.Create(ctx, artifact)
}
