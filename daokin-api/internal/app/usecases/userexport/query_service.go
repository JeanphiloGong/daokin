package userexport

import (
	"context"
	"strings"

	"daokin-api/internal/app/ports/out"
	"daokin-api/internal/domain"
)

type QueryService struct {
	artifacts   outport.ArtifactRepository
	memberships outport.MembershipRepository
	permissions outport.PermissionRepository
	transfers   outport.TransferRepository
}

func NewQueryService(
	artifacts outport.ArtifactRepository,
	memberships outport.MembershipRepository,
	permissions outport.PermissionRepository,
	transfers outport.TransferRepository,
) *QueryService {
	return &QueryService{
		artifacts:   artifacts,
		memberships: memberships,
		permissions: permissions,
		transfers:   transfers,
	}
}

func (s *QueryService) ExportByWallet(ctx context.Context, wallet string) (domain.UserExportBundle, error) {
	wallet = strings.TrimSpace(wallet)
	if wallet == "" {
		return domain.UserExportBundle{}, domain.ErrInvalidArgument
	}

	artifacts, err := s.artifacts.ListByCreator(ctx, wallet)
	if err != nil {
		return domain.UserExportBundle{}, err
	}
	memberships, err := s.memberships.ListByWalletMemberships(ctx, wallet)
	if err != nil {
		return domain.UserExportBundle{}, err
	}
	permissions, err := s.permissions.ListByWalletPermissions(ctx, wallet)
	if err != nil {
		return domain.UserExportBundle{}, err
	}
	transfers, err := s.transfers.ListByWalletTransfers(ctx, wallet)
	if err != nil {
		return domain.UserExportBundle{}, err
	}

	return domain.UserExportBundle{
		Wallet:      wallet,
		Artifacts:   artifacts,
		Memberships: memberships,
		Permissions: permissions,
		Transfers:   transfers,
	}, nil
}
