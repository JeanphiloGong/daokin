package domain

import (
	"strings"
	"time"
)

const (
	PermissionStatusActive  = "active"
	PermissionStatusRevoked = "revoked"
)

type Permission struct {
	ID            string     `json:"id"`
	ArtifactID    string     `json:"artifact_id"`
	GranterWallet string     `json:"granter_wallet"`
	GranteeWallet string     `json:"grantee_wallet"`
	Scope         string     `json:"scope"`
	Status        string     `json:"status"`
	GrantedAt     time.Time  `json:"granted_at"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	RevokeReason  string     `json:"revoke_reason,omitempty"`
}

func NewPermission(
	artifactID string,
	granterWallet string,
	granteeWallet string,
	scope string,
	grantedAt time.Time,
) (Permission, error) {
	artifactID = strings.TrimSpace(artifactID)
	granterWallet = strings.TrimSpace(granterWallet)
	granteeWallet = strings.TrimSpace(granteeWallet)
	scope = strings.TrimSpace(scope)
	if artifactID == "" || granterWallet == "" || granteeWallet == "" || scope == "" {
		return Permission{}, ErrInvalidArgument
	}
	if strings.EqualFold(granterWallet, granteeWallet) {
		return Permission{}, ErrRuleViolation
	}

	return Permission{
		ArtifactID:    artifactID,
		GranterWallet: granterWallet,
		GranteeWallet: granteeWallet,
		Scope:         scope,
		Status:        PermissionStatusActive,
		GrantedAt:     grantedAt.UTC(),
	}, nil
}
