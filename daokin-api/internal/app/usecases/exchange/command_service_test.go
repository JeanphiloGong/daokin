package exchange

import (
	"context"
	"testing"

	"daokin-api/internal/domain"
	"daokin-api/internal/infra/persistence/memory"
)

func TestRecordTransferRequiresPermissionForNonOwner(t *testing.T) {
	store := memory.NewStore()
	artifact, err := store.Create(context.Background(), domain.Artifact{
		CreatorWallet: "0xowner",
		ContentHash:   "hash",
		ContentURI:    "ipfs://cid",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	svc := NewCommandService(store, store, store)
	_, err = svc.RecordTransfer(context.Background(), artifact.ID, "0xviewer", "U", "10", "0xtx")
	if err == nil || err != domain.ErrForbidden {
		t.Fatalf("RecordTransfer() error = %v, want %v", err, domain.ErrForbidden)
	}
}
