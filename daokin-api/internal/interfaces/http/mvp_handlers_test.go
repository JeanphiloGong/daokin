package httpiface

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"daokin-api/internal/domain"
)

type fakeAuthCommands struct {
	validateErr error
}

func (f fakeAuthCommands) CreateChallenge(context.Context, string) (domain.AuthChallenge, error) {
	return domain.AuthChallenge{}, nil
}

func (f fakeAuthCommands) VerifyChallenge(context.Context, string, string, string) (domain.AuthSession, error) {
	return domain.AuthSession{}, nil
}

func (f fakeAuthCommands) ValidateSession(context.Context, string, string) error {
	return f.validateErr
}

type fakeArtifactCommands struct {
	called bool
}

func (f *fakeArtifactCommands) CreateArtifact(context.Context, string, string, string) (domain.Artifact, error) {
	f.called = true
	return domain.Artifact{ID: "art_1"}, nil
}

type fakeArtifactQueries struct{}

func (fakeArtifactQueries) GetArtifact(context.Context, string) (domain.Artifact, error) {
	return domain.Artifact{}, nil
}

type fakeMembershipCommands struct{}

func (fakeMembershipCommands) JoinDAO(context.Context, string, string) (domain.DAOMembership, error) {
	return domain.DAOMembership{}, nil
}

func (fakeMembershipCommands) LeaveDAO(context.Context, string, string, string) (domain.DAOMembership, error) {
	return domain.DAOMembership{}, nil
}

func TestCreateArtifactRequiresAuthorization(t *testing.T) {
	artifactCmd := &fakeArtifactCommands{}
	h := NewMVPHandler(
		fakeAuthCommands{},
		artifactCmd,
		fakeArtifactQueries{},
		fakeMembershipCommands{},
	)

	req := httptest.NewRequest(http.MethodPost, "/v1/artifacts", strings.NewReader(`{
		"creator_wallet":"0xabc",
		"content_hash":"hash",
		"content_uri":"ipfs://cid"
	}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateArtifact(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if artifactCmd.called {
		t.Fatalf("artifact command should not be called when unauthorized")
	}
}
