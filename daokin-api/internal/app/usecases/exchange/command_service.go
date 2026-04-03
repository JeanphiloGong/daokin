package exchange

import (
	"context"
	"strings"
	"time"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type CommandService struct {
	transfers   outport.TransferRepository
	permissions outport.PermissionRepository
	artifacts   outport.ArtifactRepository
	now         func() time.Time
}

func NewCommandService(
	transfers outport.TransferRepository,
	permissions outport.PermissionRepository,
	artifacts outport.ArtifactRepository,
) *CommandService {
	return &CommandService{
		transfers:   transfers,
		permissions: permissions,
		artifacts:   artifacts,
		now:         time.Now,
	}
}

func (s *CommandService) RecordTransfer(
	ctx context.Context,
	artifactID string,
	payerWallet string,
	token string,
	amountAtomic string,
	txHash string,
) (domain.Transfer, error) {
	artifactID = strings.TrimSpace(artifactID)
	payerWallet = strings.TrimSpace(payerWallet)

	artifact, err := s.artifacts.GetByID(ctx, artifactID)
	if err != nil {
		return domain.Transfer{}, err
	}

	if !strings.EqualFold(artifact.CreatorWallet, payerWallet) {
		ok, err := s.permissions.HasActive(ctx, artifactID, payerWallet, "view")
		if err != nil {
			return domain.Transfer{}, err
		}
		if !ok {
			return domain.Transfer{}, domain.ErrForbidden
		}
	}

	transfer, err := domain.NewTransfer(
		artifactID,
		payerWallet,
		artifact.CreatorWallet,
		token,
		amountAtomic,
		txHash,
		s.now(),
	)
	if err != nil {
		return domain.Transfer{}, err
	}
	return s.transfers.CreateTransfer(ctx, transfer)
}
