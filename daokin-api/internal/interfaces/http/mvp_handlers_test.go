package httpiface

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	artifactuc "daokin-api/internal/app/usecases/artifact"
	attributionuc "daokin-api/internal/app/usecases/attribution"
	exchangeuc "daokin-api/internal/app/usecases/exchange"
	identityuc "daokin-api/internal/app/usecases/identity"
	membershipuc "daokin-api/internal/app/usecases/membership"
	permissionuc "daokin-api/internal/app/usecases/permission"
	userexportuc "daokin-api/internal/app/usecases/userexport"
	"daokin-api/internal/domain"
	"daokin-api/internal/infra/persistence/memory"
)

const (
	ownerWallet   = "0x1111111111111111111111111111111111111111"
	granteeWallet = "0x2222222222222222222222222222222222222222"
)

type stubVerifier struct{}

func (stubVerifier) VerifyWalletSignature(context.Context, string, string, string) error {
	return nil
}

func newMVPTestHandler() http.Handler {
	store := memory.NewStore()
	mvp := NewMVPHandler(
		identityuc.NewCommandService(store, stubVerifier{}),
		artifactuc.NewCommandService(store),
		artifactuc.NewQueryService(store),
		membershipuc.NewCommandService(store),
		permissionuc.NewCommandService(store, store),
		exchangeuc.NewCommandService(store, store, store),
		attributionuc.NewQueryService(store, store),
		userexportuc.NewQueryService(store, store, store, store),
	)
	return Routes(slog.New(slog.NewTextHandler(io.Discard, nil)), mvp)
}

func performJSON(
	t *testing.T,
	handler http.Handler,
	method string,
	path string,
	token string,
	body string,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func decodeResponseJSON[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()

	var out T
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode response error = %v; body = %q", err, rr.Body.String())
	}
	return out
}

func authenticateWallet(t *testing.T, handler http.Handler, wallet string) string {
	t.Helper()

	challengeRR := performJSON(
		t,
		handler,
		http.MethodPost,
		"/v1/auth/challenge",
		"",
		`{"wallet":"`+wallet+`"}`,
	)
	if challengeRR.Code != http.StatusOK {
		t.Fatalf("challenge status = %d, want %d; body = %s", challengeRR.Code, http.StatusOK, challengeRR.Body.String())
	}
	challenge := decodeResponseJSON[domain.AuthChallenge](t, challengeRR)

	verifyRR := performJSON(
		t,
		handler,
		http.MethodPost,
		"/v1/auth/verify",
		"",
		`{"wallet":"`+wallet+`","nonce":"`+challenge.Nonce+`","signature":"0xsig"}`,
	)
	if verifyRR.Code != http.StatusOK {
		t.Fatalf("verify status = %d, want %d; body = %s", verifyRR.Code, http.StatusOK, verifyRR.Body.String())
	}
	session := decodeResponseJSON[domain.AuthSession](t, verifyRR)
	if session.AccessToken == "" {
		t.Fatalf("verify returned empty access token")
	}
	return session.AccessToken
}

func TestCreateArtifactRequiresAuthorization(t *testing.T) {
	handler := newMVPTestHandler()

	rr := performJSON(t, handler, http.MethodPost, "/v1/artifacts", "", `{
		"creator_wallet":"`+ownerWallet+`",
		"content_hash":"hash",
		"content_uri":"ipfs://cid"
	}`)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestMVPHandlerCoreLoop(t *testing.T) {
	handler := newMVPTestHandler()
	ownerToken := authenticateWallet(t, handler, ownerWallet)
	granteeToken := authenticateWallet(t, handler, granteeWallet)

	createRR := performJSON(t, handler, http.MethodPost, "/v1/artifacts", ownerToken, `{
		"creator_wallet":"`+ownerWallet+`",
		"content_hash":"0xabc",
		"content_uri":"daokin://artifact/test"
	}`)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("create artifact status = %d, want %d; body = %s", createRR.Code, http.StatusCreated, createRR.Body.String())
	}
	artifact := decodeResponseJSON[domain.Artifact](t, createRR)
	if artifact.ID == "" || artifact.CreatorWallet != ownerWallet {
		t.Fatalf("artifact = %+v, want id and owner wallet", artifact)
	}

	joinRR := performJSON(t, handler, http.MethodPost, "/v1/daos/dao-builders/join", ownerToken, `{
		"wallet":"`+ownerWallet+`"
	}`)
	if joinRR.Code != http.StatusOK {
		t.Fatalf("join status = %d, want %d; body = %s", joinRR.Code, http.StatusOK, joinRR.Body.String())
	}
	membership := decodeResponseJSON[domain.DAOMembership](t, joinRR)
	if membership.DAOID != "dao-builders" || membership.Wallet != ownerWallet || membership.LeftAt != nil {
		t.Fatalf("membership = %+v, want active dao-builders membership", membership)
	}

	deniedTransferRR := performJSON(t, handler, http.MethodPost, "/v1/transfers", granteeToken, `{
		"artifact_id":"`+artifact.ID+`",
		"payer_wallet":"`+granteeWallet+`",
		"token":"USDC",
		"amount_atomic":"1000000",
		"tx_hash":"0xdenied"
	}`)
	if deniedTransferRR.Code != http.StatusForbidden {
		t.Fatalf("denied transfer status = %d, want %d; body = %s", deniedTransferRR.Code, http.StatusForbidden, deniedTransferRR.Body.String())
	}

	grantRR := performJSON(t, handler, http.MethodPost, "/v1/permissions", ownerToken, `{
		"artifact_id":"`+artifact.ID+`",
		"granter_wallet":"`+ownerWallet+`",
		"grantee_wallet":"`+granteeWallet+`",
		"scope":"view"
	}`)
	if grantRR.Code != http.StatusCreated {
		t.Fatalf("grant status = %d, want %d; body = %s", grantRR.Code, http.StatusCreated, grantRR.Body.String())
	}
	permission := decodeResponseJSON[domain.Permission](t, grantRR)
	if permission.Status != domain.PermissionStatusActive || permission.GranteeWallet != granteeWallet {
		t.Fatalf("permission = %+v, want active grantee permission", permission)
	}

	transferRR := performJSON(t, handler, http.MethodPost, "/v1/transfers", granteeToken, `{
		"artifact_id":"`+artifact.ID+`",
		"payer_wallet":"`+granteeWallet+`",
		"token":"USDC",
		"amount_atomic":"1000000",
		"tx_hash":"0xallowed"
	}`)
	if transferRR.Code != http.StatusCreated {
		t.Fatalf("transfer status = %d, want %d; body = %s", transferRR.Code, http.StatusCreated, transferRR.Body.String())
	}
	transfer := decodeResponseJSON[domain.Transfer](t, transferRR)
	if transfer.RecipientWallet != ownerWallet || transfer.PayerWallet != granteeWallet {
		t.Fatalf("transfer = %+v, want payer grantee and recipient owner", transfer)
	}

	attributionRR := performJSON(
		t,
		handler,
		http.MethodGet,
		"/v1/artifacts/"+artifact.ID+"/attribution",
		"",
		"",
	)
	if attributionRR.Code != http.StatusOK {
		t.Fatalf("attribution status = %d, want %d; body = %s", attributionRR.Code, http.StatusOK, attributionRR.Body.String())
	}
	trail := decodeResponseJSON[domain.AttributionTrail](t, attributionRR)
	if trail.TotalRecords != 2 || len(trail.Permissions) != 1 || len(trail.Transfers) != 1 {
		t.Fatalf("trail = %+v, want one permission and one transfer", trail)
	}

	exportRR := performJSON(t, handler, http.MethodGet, "/v1/users/"+ownerWallet+"/export", ownerToken, "")
	if exportRR.Code != http.StatusOK {
		t.Fatalf("export status = %d, want %d; body = %s", exportRR.Code, http.StatusOK, exportRR.Body.String())
	}
	bundle := decodeResponseJSON[domain.UserExportBundle](t, exportRR)
	if len(bundle.Artifacts) != 1 || len(bundle.Memberships) != 1 || len(bundle.Permissions) != 1 || len(bundle.Transfers) != 1 {
		t.Fatalf("bundle = %+v, want owner artifact, membership, permission, and received transfer", bundle)
	}

	revokeRR := performJSON(t, handler, http.MethodPost, "/v1/permissions/"+permission.ID+"/revoke", ownerToken, `{
		"granter_wallet":"`+ownerWallet+`",
		"reason":"test revoke"
	}`)
	if revokeRR.Code != http.StatusOK {
		t.Fatalf("revoke status = %d, want %d; body = %s", revokeRR.Code, http.StatusOK, revokeRR.Body.String())
	}
	revoked := decodeResponseJSON[domain.Permission](t, revokeRR)
	if revoked.Status != domain.PermissionStatusRevoked || revoked.RevokedAt == nil {
		t.Fatalf("revoked permission = %+v, want revoked status and timestamp", revoked)
	}

	leaveRR := performJSON(t, handler, http.MethodPost, "/v1/daos/dao-builders/leave", ownerToken, `{
		"wallet":"`+ownerWallet+`",
		"reason":"test leave"
	}`)
	if leaveRR.Code != http.StatusOK {
		t.Fatalf("leave status = %d, want %d; body = %s", leaveRR.Code, http.StatusOK, leaveRR.Body.String())
	}
	left := decodeResponseJSON[domain.DAOMembership](t, leaveRR)
	if left.LeftAt == nil || left.LeaveReason != "test leave" {
		t.Fatalf("left membership = %+v, want left_at and reason", left)
	}
}
