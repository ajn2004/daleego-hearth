// internal/httpapi/handlers/mobile/handler.go
package mobile

import (
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/ajn2004/daleego-hearth/backend/internal/middleware"
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

func (h *Handler) Routes(queries *db.Queries) chi.Router {
	r := chi.NewRouter()

	h.registerPairingRoutes(r)
	r.Group(func(r chi.Router) {
		r.Use(middleware.DeviceAuth(queries))
		h.registerLocationRoutes(r)
	})
	return r
}
