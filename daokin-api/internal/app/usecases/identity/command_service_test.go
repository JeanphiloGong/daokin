package identity

import (
	"context"
	"testing"

	"daokin-api/internal/domain"
	"daokin-api/internal/infra/persistence/memory"
)

type stubVerifier struct {
	verify func(wallet, message, signature string) error
}

func (s stubVerifier) VerifyWalletSignature(_ context.Context, wallet, message, signature string) error {
	if s.verify == nil {
		return nil
	}
	return s.verify(wallet, message, signature)
}

func TestVerifyChallengeRejectsInvalidSignature(t *testing.T) {
	store := memory.NewStore()
	svc := NewCommandService(store, stubVerifier{
		verify: func(_, _, _ string) error { return domain.ErrUnauthorized },
	})

	ctx := context.Background()
	challenge, err := svc.CreateChallenge(ctx, "0xabc")
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	_, err = svc.VerifyChallenge(ctx, challenge.Wallet, challenge.Nonce, "bad-signature")
	if err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("VerifyChallenge() error = %v, want %v", err, domain.ErrUnauthorized)
	}
}

func TestVerifyChallengeConsumesChallengeOnce(t *testing.T) {
	store := memory.NewStore()
	svc := NewCommandService(store, stubVerifier{})

	ctx := context.Background()
	challenge, err := svc.CreateChallenge(ctx, "0xabc")
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	session, err := svc.VerifyChallenge(ctx, challenge.Wallet, challenge.Nonce, "ok")
	if err != nil {
		t.Fatalf("VerifyChallenge() first error = %v", err)
	}
	if session.AccessToken == "" {
		t.Fatalf("VerifyChallenge() empty token")
	}

	_, err = svc.VerifyChallenge(ctx, challenge.Wallet, challenge.Nonce, "ok")
	if err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("VerifyChallenge() second error = %v, want %v", err, domain.ErrUnauthorized)
	}
}

func TestValidateSessionRejectsTokenMismatch(t *testing.T) {
	store := memory.NewStore()
	svc := NewCommandService(store, stubVerifier{})

	ctx := context.Background()
	challenge, err := svc.CreateChallenge(ctx, "0xabc")
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}
	session, err := svc.VerifyChallenge(ctx, challenge.Wallet, challenge.Nonce, "ok")
	if err != nil {
		t.Fatalf("VerifyChallenge() error = %v", err)
	}

	err = svc.ValidateSession(ctx, challenge.Wallet, session.AccessToken+"x")
	if err == nil || err != domain.ErrUnauthorized {
		t.Fatalf("ValidateSession() error = %v, want %v", err, domain.ErrUnauthorized)
	}
}
