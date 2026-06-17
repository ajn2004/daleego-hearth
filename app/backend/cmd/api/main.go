package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/ajn2004/daleego-hearth/backend/internal/config"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}
	defer pool.Close()
	addr := ":" + cfg.Port
	slog.Info("Starting api server", "addr", addr)
	if err := http.ListenAndServe(addr, httpapi.NewRouter(pool, cfg.AdminAPIKey)); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}
