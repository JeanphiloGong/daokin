package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

var (
	_ outport.AuthRepository       = (*Store)(nil)
	_ outport.ArtifactRepository   = (*Store)(nil)
	_ outport.MembershipRepository = (*Store)(nil)
	_ outport.PermissionRepository = (*Store)(nil)
	_ outport.TransferRepository   = (*Store)(nil)
)

type Store struct {
	mu sync.RWMutex

	challenges    map[string]domain.AuthChallenge
	sessions      map[string]domain.AuthSession
	artifacts     map[string]domain.Artifact
	memberships   map[string]domain.DAOMembership
	permissions   map[string]domain.Permission
	transfers     map[string]domain.Transfer
	artifactSeq   int
	permissionSeq int
	transferSeq   int
}

func NewStore() *Store {
	return &Store{
		challenges:  make(map[string]domain.AuthChallenge),
		sessions:    make(map[string]domain.AuthSession),
		artifacts:   make(map[string]domain.Artifact),
		memberships: make(map[string]domain.DAOMembership),
		permissions: make(map[string]domain.Permission),
		transfers:   make(map[string]domain.Transfer),
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

func (s *Store) ConsumeChallenge(_ context.Context, wallet string) (domain.AuthChallenge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := normalizeWallet(wallet)
	challenge, ok := s.challenges[key]
	if !ok {
		return domain.AuthChallenge{}, domain.ErrNotFound
	}
	delete(s.challenges, key)
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

func (s *Store) GetSession(_ context.Context, wallet string) (domain.AuthSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[normalizeWallet(wallet)]
	if !ok {
		return domain.AuthSession{}, domain.ErrNotFound
	}
	return session, nil
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

func (s *Store) ListByCreator(_ context.Context, wallet string) ([]domain.Artifact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedWallet := normalizeWallet(wallet)
	artifacts := make([]domain.Artifact, 0)
	for _, artifact := range s.artifacts {
		if normalizeWallet(artifact.CreatorWallet) == normalizedWallet {
			artifacts = append(artifacts, artifact)
		}
	}
	return artifacts, nil
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

func (s *Store) ListByWalletMemberships(_ context.Context, wallet string) ([]domain.DAOMembership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedWallet := normalizeWallet(wallet)
	memberships := make([]domain.DAOMembership, 0)
	for _, membership := range s.memberships {
		if normalizeWallet(membership.Wallet) == normalizedWallet {
			memberships = append(memberships, membership)
		}
	}
	return memberships, nil
}

func (s *Store) Grant(_ context.Context, permission domain.Permission) (domain.Permission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, current := range s.permissions {
		if current.ArtifactID == permission.ArtifactID &&
			normalizeWallet(current.GranteeWallet) == normalizeWallet(permission.GranteeWallet) &&
			strings.EqualFold(current.Scope, permission.Scope) &&
			current.Status == domain.PermissionStatusActive {
			return current, nil
		}
	}

	s.permissionSeq++
	permission.ID = fmt.Sprintf("perm_%06d", s.permissionSeq)
	s.permissions[permission.ID] = permission
	return permission, nil
}

func (s *Store) Revoke(
	_ context.Context,
	permissionID, granterWallet string,
	revokedAt time.Time,
	reason string,
) (domain.Permission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.permissions[strings.TrimSpace(permissionID)]
	if !ok {
		return domain.Permission{}, domain.ErrNotFound
	}
	if !strings.EqualFold(current.GranterWallet, strings.TrimSpace(granterWallet)) {
		return domain.Permission{}, domain.ErrForbidden
	}
	if current.Status != domain.PermissionStatusActive {
		return domain.Permission{}, domain.ErrConflict
	}

	revokedAtCopy := revokedAt.UTC()
	current.Status = domain.PermissionStatusRevoked
	current.RevokedAt = &revokedAtCopy
	current.RevokeReason = strings.TrimSpace(reason)
	s.permissions[current.ID] = current
	return current, nil
}

func (s *Store) HasActive(_ context.Context, artifactID, granteeWallet, scope string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	artifactID = strings.TrimSpace(artifactID)
	granteeWallet = normalizeWallet(granteeWallet)
	scope = strings.TrimSpace(scope)
	for _, permission := range s.permissions {
		if permission.ArtifactID == artifactID &&
			normalizeWallet(permission.GranteeWallet) == granteeWallet &&
			strings.EqualFold(permission.Scope, scope) &&
			permission.Status == domain.PermissionStatusActive {
			return true, nil
		}
	}
	return false, nil
}

func (s *Store) ListPermissionsByArtifact(_ context.Context, artifactID string) ([]domain.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	artifactID = strings.TrimSpace(artifactID)
	permissions := make([]domain.Permission, 0)
	for _, permission := range s.permissions {
		if permission.ArtifactID == artifactID {
			permissions = append(permissions, permission)
		}
	}
	return permissions, nil
}

func (s *Store) ListByWalletPermissions(_ context.Context, wallet string) ([]domain.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedWallet := normalizeWallet(wallet)
	permissions := make([]domain.Permission, 0)
	for _, permission := range s.permissions {
		if normalizeWallet(permission.GranteeWallet) == normalizedWallet ||
			normalizeWallet(permission.GranterWallet) == normalizedWallet {
			permissions = append(permissions, permission)
		}
	}
	return permissions, nil
}

func (s *Store) CreateTransfer(_ context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transferSeq++
	transfer.ID = fmt.Sprintf("tx_%06d", s.transferSeq)
	s.transfers[transfer.ID] = transfer
	return transfer, nil
}

func (s *Store) ListTransfersByArtifact(_ context.Context, artifactID string) ([]domain.Transfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	artifactID = strings.TrimSpace(artifactID)
	transfers := make([]domain.Transfer, 0)
	for _, transfer := range s.transfers {
		if transfer.ArtifactID == artifactID {
			transfers = append(transfers, transfer)
		}
	}
	return transfers, nil
}

func (s *Store) ListByWalletTransfers(_ context.Context, wallet string) ([]domain.Transfer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedWallet := normalizeWallet(wallet)
	transfers := make([]domain.Transfer, 0)
	for _, transfer := range s.transfers {
		if normalizeWallet(transfer.PayerWallet) == normalizedWallet ||
			normalizeWallet(transfer.RecipientWallet) == normalizedWallet {
			transfers = append(transfers, transfer)
		}
	}
	return transfers, nil
}

func normalizeWallet(wallet string) string {
	return strings.ToLower(strings.TrimSpace(wallet))
}

func membershipKey(daoID, wallet string) string {
	return strings.TrimSpace(daoID) + "|" + normalizeWallet(wallet)
}
