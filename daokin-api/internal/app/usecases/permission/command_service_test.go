package permission

import (
	"context"
	"testing"

	"daokin-api/internal/domain"
	"daokin-api/internal/infra/persistence/memory"
)

func TestGrantPermissionOnlyOwnerCanGrant(t *testing.T) {
	store := memory.NewStore()
	artifact, err := store.Create(context.Background(), domain.Artifact{
		CreatorWallet: "0xowner",
		ContentHash:   "hash",
		ContentURI:    "ipfs://cid",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	svc := NewCommandService(store, store)
	_, err = svc.GrantPermission(context.Background(), artifact.ID, "0xother", "0xviewer", "view")
	if err == nil || err != domain.ErrForbidden {
		t.Fatalf("GrantPermission() error = %v, want %v", err, domain.ErrForbidden)
	}
}
