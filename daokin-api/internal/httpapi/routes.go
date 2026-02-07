package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"daokin-api/internal/httpapi/handlers"
	appmw "daokin-api/internal/httpapi/middleware"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
)

func Routes(logger *slog.Logger, mvp *MVPHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimid.RequestID)
	r.Use(chimid.RealIP)
	r.Use(appmw.RequestLogger(logger))
	r.Use(chimid.Recoverer)
	r.Use(chimid.Timeout(15 * time.Second))

	r.Get("/healthz", handlers.Health)
	r.Get("/readyz", handlers.Ready)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
		})

		r.Post("/auth/challenge", mvp.AuthChallenge)
		r.Post("/auth/verify", mvp.AuthVerify)
		r.Post("/artifacts", mvp.CreateArtifact)
		r.Get("/artifacts/{id}", mvp.GetArtifact)
		r.Post("/daos/{id}/join", mvp.JoinDAO)
		r.Post("/daos/{id}/leave", mvp.LeaveDAO)
	})

	return r
}
