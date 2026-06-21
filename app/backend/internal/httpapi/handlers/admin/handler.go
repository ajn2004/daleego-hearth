// internal/httpapi/handlers/admin/handler.go
package admin

import (
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	h.registerPeopleRoutes(r)
	h.registerDeviceRoutes(r)
	h.registerPairingRoutes(r)
	h.registerLocationRoutes(r)
	return r
}
