package httpiface

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"daokin-api/internal/app/ports/in"
	"daokin-api/internal/domain"

	"github.com/go-chi/chi/v5"
)

type MVPHandler struct {
	authCommands       inport.AuthCommands
	artifactCommands   inport.ArtifactCommands
	artifactQueries    inport.ArtifactQueries
	membershipCommands inport.MembershipCommands
}

func NewMVPHandler(
	authCommands inport.AuthCommands,
	artifactCommands inport.ArtifactCommands,
	artifactQueries inport.ArtifactQueries,
	membershipCommands inport.MembershipCommands,
) *MVPHandler {
	return &MVPHandler{
		authCommands:       authCommands,
		artifactCommands:   artifactCommands,
		artifactQueries:    artifactQueries,
		membershipCommands: membershipCommands,
	}
}

func (h *MVPHandler) AuthChallenge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Wallet string `json:"wallet"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	challenge, err := h.authCommands.CreateChallenge(r.Context(), req.Wallet)
	if err != nil {
		writeDomainError(w, err, "failed to create challenge")
		return
	}

	writeJSON(w, http.StatusOK, challenge)
}

func (h *MVPHandler) AuthVerify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Wallet    string `json:"wallet"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := h.authCommands.VerifyChallenge(r.Context(), req.Wallet, req.Nonce, req.Signature)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, "challenge verification failed")
			return
		}
		writeDomainError(w, err, "failed to verify challenge")
		return
	}

	writeJSON(w, http.StatusOK, session)
}

func (h *MVPHandler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CreatorWallet string `json:"creator_wallet"`
		ContentHash   string `json:"content_hash"`
		ContentURI    string `json:"content_uri"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.authorizeWallet(r, req.CreatorWallet); err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	artifact, err := h.artifactCommands.CreateArtifact(
		r.Context(),
		req.CreatorWallet,
		req.ContentHash,
		req.ContentURI,
	)
	if err != nil {
		writeDomainError(w, err, "failed to create artifact")
		return
	}

	writeJSON(w, http.StatusCreated, artifact)
}

func (h *MVPHandler) GetArtifact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	artifact, err := h.artifactQueries.GetArtifact(r.Context(), id)
	if err != nil {
		writeDomainError(w, err, "failed to get artifact")
		return
	}

	writeJSON(w, http.StatusOK, artifact)
}

func (h *MVPHandler) JoinDAO(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Wallet string `json:"wallet"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.authorizeWallet(r, req.Wallet); err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	membership, err := h.membershipCommands.JoinDAO(r.Context(), chi.URLParam(r, "id"), req.Wallet)
	if err != nil {
		writeDomainError(w, err, "failed to join dao")
		return
	}

	writeJSON(w, http.StatusOK, membership)
}

func (h *MVPHandler) LeaveDAO(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Wallet string `json:"wallet"`
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.authorizeWallet(r, req.Wallet); err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	membership, err := h.membershipCommands.LeaveDAO(r.Context(), chi.URLParam(r, "id"), req.Wallet, req.Reason)
	if err != nil {
		writeDomainError(w, err, "failed to leave dao")
		return
	}

	writeJSON(w, http.StatusOK, membership)
}

func (h *MVPHandler) authorizeWallet(r *http.Request, wallet string) error {
	token, err := bearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return domain.ErrUnauthorized
	}
	return h.authCommands.ValidateSession(r.Context(), wallet, token)
}

func bearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", domain.ErrUnauthorized
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", domain.ErrUnauthorized
	}
	return parts[1], nil
}

func writeDomainError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		writeError(w, http.StatusBadRequest, "invalid argument")
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, domain.ErrConflict):
		writeError(w, http.StatusConflict, "state conflict")
	case errors.Is(err, domain.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	default:
		writeError(w, http.StatusInternalServerError, fallback)
	}
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}

	var trailing struct{}
	if err := dec.Decode(&trailing); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
