package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"daokin-api/internal/domain"
	"daokin-api/internal/repository"
)

var (
	_ repository.AuthRepository       = (*Store)(nil)
	_ repository.ArtifactRepository   = (*Store)(nil)
	_ repository.MembershipRepository = (*Store)(nil)
)

type Store struct {
	mu sync.RWMutex

	challenges  map[string]domain.AuthChallenge
	sessions    map[string]domain.AuthSession
	artifacts   map[string]domain.Artifact
	memberships map[string]domain.DAOMembership
	artifactSeq int
}

func NewStore() *Store {
	return &Store{
		challenges:  make(map[string]domain.AuthChallenge),
		sessions:    make(map[string]domain.AuthSession),
		artifacts:   make(map[string]domain.Artifact),
		memberships: make(map[string]domain.DAOMembership),
	}
}

func (s *Store) SaveChallenge(_ context.Context, challenge domain.AuthChallenge) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.challenges[normalizeWallet(challenge.Wallet)] = challenge
	return nil
}

func (s *Store) GetChallenge(_ context.Context, wallet string) (domain.AuthChallenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	challenge, ok := s.challenges[normalizeWallet(wallet)]
	if !ok {
		return domain.AuthChallenge{}, domain.ErrNotFound
	}
	return challenge, nil
}

func (s *Store) DeleteChallenge(_ context.Context, wallet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.challenges, normalizeWallet(wallet))
	return nil
}

func (s *Store) SaveSession(_ context.Context, session domain.AuthSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[normalizeWallet(session.Wallet)] = session
	return nil
}

func (s *Store) Create(_ context.Context, artifact domain.Artifact) (domain.Artifact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.artifactSeq++
	artifact.ID = fmt.Sprintf("art_%06d", s.artifactSeq)
	s.artifacts[artifact.ID] = artifact
	return artifact, nil
}

func (s *Store) GetByID(_ context.Context, id string) (domain.Artifact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	artifact, ok := s.artifacts[id]
	if !ok {
		return domain.Artifact{}, domain.ErrNotFound
	}
	return artifact, nil
}

func (s *Store) Join(_ context.Context, membership domain.DAOMembership) (domain.DAOMembership, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := membershipKey(membership.DAOID, membership.Wallet)
	current, ok := s.memberships[key]
	if !ok {
		s.memberships[key] = membership
		return membership, nil
	}

	if current.LeftAt == nil {
		return current, nil
	}

	current.JoinedAt = membership.JoinedAt
	current.LeftAt = nil
	current.LeaveReason = ""
	s.memberships[key] = current
	return current, nil
}

func (s *Store) Leave(_ context.Context, daoID, wallet string, leftAt time.Time, reason string) (domain.DAOMembership, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := membershipKey(daoID, wallet)
	current, ok := s.memberships[key]
	if !ok {
		return domain.DAOMembership{}, domain.ErrNotFound
	}
	if current.LeftAt != nil {
		return domain.DAOMembership{}, domain.ErrConflict
	}

	leftAtCopy := leftAt
	current.LeftAt = &leftAtCopy
	current.LeaveReason = strings.TrimSpace(reason)
	s.memberships[key] = current
	return current, nil
}

func normalizeWallet(wallet string) string {
	return strings.ToLower(strings.TrimSpace(wallet))
}

func membershipKey(daoID, wallet string) string {
	return strings.TrimSpace(daoID) + "|" + normalizeWallet(wallet)
}
