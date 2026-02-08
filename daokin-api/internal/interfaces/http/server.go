package httpiface

import (
	"log/slog"
	"net/http"

	artifactuc "daokin-api/internal/app/usecases/artifact"
	identityuc "daokin-api/internal/app/usecases/identity"
	membershipuc "daokin-api/internal/app/usecases/membership"
	"daokin-api/internal/config"
	"daokin-api/internal/infra/persistence/memory"
)

func NewServer(cfg config.Config, logger *slog.Logger) *http.Server {
	store := memory.NewStore()
	mvp := NewMVPHandler(
		identityuc.NewCommandService(store),
		artifactuc.NewCommandService(store),
		artifactuc.NewQueryService(store),
		membershipuc.NewCommandService(store),
	)

	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      Routes(logger, mvp),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
