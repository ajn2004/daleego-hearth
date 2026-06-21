// internal/httpapi/handlers/mobile/handler.go
package mobile

import (
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	queries *db.Queries
	dbPool  *pgxpool.Pool
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	h.registerLocationRoutes(r)
	h.registerPairingRoutes(r)
	return r
}
