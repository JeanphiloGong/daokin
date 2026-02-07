package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"daokin-api/internal/domain"
	"daokin-api/internal/repository"

	"github.com/go-chi/chi/v5"
)

type MVPHandler struct {
	authRepo       repository.AuthRepository
	artifactRepo   repository.ArtifactRepository
	membershipRepo repository.MembershipRepository
	now            func() time.Time
}

func NewMVPHandler(
	authRepo repository.AuthRepository,
	artifactRepo repository.ArtifactRepository,
	membershipRepo repository.MembershipRepository,
) *MVPHandler {
	return &MVPHandler{
		authRepo:       authRepo,
		artifactRepo:   artifactRepo,
		membershipRepo: membershipRepo,
		now:            time.Now,
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

	wallet := strings.TrimSpace(req.Wallet)
	if wallet == "" {
		writeError(w, http.StatusBadRequest, "wallet is required")
		return
	}

	now := h.now().UTC()
	nonce, err := randomHex(16)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create challenge")
		return
	}

	challenge := domain.AuthChallenge{
		Wallet:    wallet,
		Nonce:     nonce,
		Message:   fmt.Sprintf("Sign this DaoKin challenge nonce: %s", nonce),
		CreatedAt: now,
		ExpiresAt: now.Add(5 * time.Minute),
	}
	if err := h.authRepo.SaveChallenge(r.Context(), challenge); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save challenge")
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

	wallet := strings.TrimSpace(req.Wallet)
	nonce := strings.TrimSpace(req.Nonce)
	signature := strings.TrimSpace(req.Signature)
	if wallet == "" || nonce == "" || signature == "" {
		writeError(w, http.StatusBadRequest, "wallet, nonce and signature are required")
		return
	}

	challenge, err := h.authRepo.GetChallenge(r.Context(), wallet)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "challenge not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load challenge")
		return
	}

	now := h.now().UTC()
	if challenge.Nonce != nonce || now.After(challenge.ExpiresAt) {
		writeError(w, http.StatusUnauthorized, "challenge verification failed")
		return
	}

	token, err := randomHex(24)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	session := domain.AuthSession{
		Wallet:      wallet,
		AccessToken: token,
		VerifiedAt:  now,
	}
	if err := h.authRepo.SaveSession(r.Context(), session); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save session")
		return
	}

	_ = h.authRepo.DeleteChallenge(r.Context(), wallet)

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

	creatorWallet := strings.TrimSpace(req.CreatorWallet)
	contentHash := strings.TrimSpace(req.ContentHash)
	contentURI := strings.TrimSpace(req.ContentURI)
	if creatorWallet == "" || contentHash == "" || contentURI == "" {
		writeError(w, http.StatusBadRequest, "creator_wallet, content_hash and content_uri are required")
		return
	}

	artifact, err := h.artifactRepo.Create(r.Context(), domain.Artifact{
		CreatorWallet: creatorWallet,
		ContentHash:   contentHash,
		ContentURI:    contentURI,
		CreatedAt:     h.now().UTC(),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create artifact")
		return
	}

	writeJSON(w, http.StatusCreated, artifact)
}

func (h *MVPHandler) GetArtifact(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "artifact id is required")
		return
	}

	artifact, err := h.artifactRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "artifact not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get artifact")
		return
	}

	writeJSON(w, http.StatusOK, artifact)
}

func (h *MVPHandler) JoinDAO(w http.ResponseWriter, r *http.Request) {
	daoID := strings.TrimSpace(chi.URLParam(r, "id"))
	if daoID == "" {
		writeError(w, http.StatusBadRequest, "dao id is required")
		return
	}

	var req struct {
		Wallet string `json:"wallet"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wallet := strings.TrimSpace(req.Wallet)
	if wallet == "" {
		writeError(w, http.StatusBadRequest, "wallet is required")
		return
	}

	membership, err := h.membershipRepo.Join(r.Context(), domain.DAOMembership{
		DAOID:    daoID,
		Wallet:   wallet,
		JoinedAt: h.now().UTC(),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to join dao")
		return
	}

	writeJSON(w, http.StatusOK, membership)
}

func (h *MVPHandler) LeaveDAO(w http.ResponseWriter, r *http.Request) {
	daoID := strings.TrimSpace(chi.URLParam(r, "id"))
	if daoID == "" {
		writeError(w, http.StatusBadRequest, "dao id is required")
		return
	}

	var req struct {
		Wallet string `json:"wallet"`
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wallet := strings.TrimSpace(req.Wallet)
	if wallet == "" {
		writeError(w, http.StatusBadRequest, "wallet is required")
		return
	}

	membership, err := h.membershipRepo.Leave(r.Context(), daoID, wallet, h.now().UTC(), req.Reason)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "membership not found")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			writeError(w, http.StatusConflict, "membership already left")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to leave dao")
		return
	}

	writeJSON(w, http.StatusOK, membership)
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

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
