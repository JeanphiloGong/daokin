package membership

import (
	"context"
	"strings"
	"time"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type CommandService struct {
	repo outport.MembershipRepository
	now  func() time.Time
}

func NewCommandService(repo outport.MembershipRepository) *CommandService {
	return &CommandService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *CommandService) JoinDAO(ctx context.Context, daoID string, wallet string) (domain.DAOMembership, error) {
	membership, err := domain.NewDAOMembership(daoID, wallet, s.now())
	if err != nil {
		return domain.DAOMembership{}, err
	}
	return s.repo.Join(ctx, membership)
}

func (s *CommandService) LeaveDAO(
	ctx context.Context,
	daoID string,
	wallet string,
	reason string,
) (domain.DAOMembership, error) {
	daoID = strings.TrimSpace(daoID)
	wallet = strings.TrimSpace(wallet)
	if daoID == "" || wallet == "" {
		return domain.DAOMembership{}, domain.ErrInvalidArgument
	}
	return s.repo.Leave(ctx, daoID, wallet, s.now().UTC(), reason)
}
