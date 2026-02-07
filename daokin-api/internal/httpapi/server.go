package httpapi

import (
	"log/slog"
	"net/http"

	"daokin-api/internal/config"
	"daokin-api/internal/repository/memory"
)

func NewServer(cfg config.Config, logger *slog.Logger) *http.Server {
	store := memory.NewStore()
	mvp := NewMVPHandler(store, store, store)

	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      Routes(logger, mvp),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
