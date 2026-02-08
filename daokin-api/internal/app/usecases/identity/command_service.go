package identity

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type CommandService struct {
	repo outport.AuthRepository
	now  func() time.Time
}

func NewCommandService(repo outport.AuthRepository) *CommandService {
	return &CommandService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *CommandService) CreateChallenge(ctx context.Context, wallet string) (domain.AuthChallenge, error) {
	nonce, err := randomHex(16)
	if err != nil {
		return domain.AuthChallenge{}, err
	}

	challenge, err := domain.NewAuthChallenge(
		wallet,
		nonce,
		fmt.Sprintf("Sign this DaoKin challenge nonce: %s", nonce),
		s.now(),
		5*time.Minute,
	)
	if err != nil {
		return domain.AuthChallenge{}, err
	}

	if err := s.repo.SaveChallenge(ctx, challenge); err != nil {
		return domain.AuthChallenge{}, err
	}
	return challenge, nil
}

func (s *CommandService) VerifyChallenge(
	ctx context.Context,
	wallet string,
	nonce string,
	signature string,
) (domain.AuthSession, error) {
	wallet = strings.TrimSpace(wallet)
	nonce = strings.TrimSpace(nonce)
	signature = strings.TrimSpace(signature)
	if wallet == "" || nonce == "" || signature == "" {
		return domain.AuthSession{}, domain.ErrInvalidArgument
	}

	challenge, err := s.repo.GetChallenge(ctx, wallet)
	if err != nil {
		return domain.AuthSession{}, err
	}

	if err := challenge.Verify(nonce, s.now()); err != nil {
		if errors.Is(err, domain.ErrChallengeInvalid) {
			return domain.AuthSession{}, domain.ErrUnauthorized
		}
		return domain.AuthSession{}, err
	}

	token, err := randomHex(24)
	if err != nil {
		return domain.AuthSession{}, err
	}

	session, err := domain.NewAuthSession(wallet, token, s.now())
	if err != nil {
		return domain.AuthSession{}, err
	}

	if err := s.repo.SaveSession(ctx, session); err != nil {
		return domain.AuthSession{}, err
	}
	_ = s.repo.DeleteChallenge(ctx, wallet)
	return session, nil
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
