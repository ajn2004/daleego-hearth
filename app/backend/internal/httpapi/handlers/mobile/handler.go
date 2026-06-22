// internal/httpapi/handlers/mobile/handler.go
package mobile

import (
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	queries           *db.Queries
	dbPool            *pgxpool.Pool
	pairingCodeSecret string
}

func NewHandler(queries *db.Queries, dbPool *pgxpool.Pool, PairingCodeSecret string) *Handler {
	return &Handler{
		queries:           queries,
		dbPool:            dbPool,
		pairingCodeSecret: PairingCodeSecret,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	h.registerLocationRoutes(r)
	h.registerPairingRoutes(r)
	return r
}
