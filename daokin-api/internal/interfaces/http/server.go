package httpiface

import (
	"log/slog"
	"net/http"

	artifactuc "daokin-api/internal/app/usecases/artifact"
	attributionuc "daokin-api/internal/app/usecases/attribution"
	exchangeuc "daokin-api/internal/app/usecases/exchange"
	identityuc "daokin-api/internal/app/usecases/identity"
	membershipuc "daokin-api/internal/app/usecases/membership"
	permissionuc "daokin-api/internal/app/usecases/permission"
	userexportuc "daokin-api/internal/app/usecases/userexport"
	"daokin-api/internal/config"
	"daokin-api/internal/infra/crypto/ethsign"
	"daokin-api/internal/infra/persistence/memory"
)

func NewServer(cfg config.Config, logger *slog.Logger) *http.Server {
	store := memory.NewStore()
	mvp := NewMVPHandler(
		identityuc.NewCommandService(store, ethsign.NewPersonalSignVerifier()),
		artifactuc.NewCommandService(store),
		artifactuc.NewQueryService(store),
		membershipuc.NewCommandService(store),
		permissionuc.NewCommandService(store, store),
		exchangeuc.NewCommandService(store, store, store),
		attributionuc.NewQueryService(store, store),
		userexportuc.NewQueryService(store, store, store, store),
	)

	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      Routes(logger, mvp),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
