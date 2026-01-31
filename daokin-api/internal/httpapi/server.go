package httpapi

import (
	"log/slog"
	"net/http"

	"daokin-api/internal/config"
)

func NewServer(cfg config.Config, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      Routes(logger),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
