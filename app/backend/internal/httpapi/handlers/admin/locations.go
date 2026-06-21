package admin

import (
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	httputil "github.com/ajn2004/daleego-hearth/backend/internal/httpapi/utils"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) registerLocationRoutes(r chi.Router) {
	r.Get("/locations/location/{location_id}", h.GetLocation)
	r.Get("/locations/people/{person_id}", h.GetLocationsByPerson)
}

func (h *Handler) GetLocation(w http.ResponseWriter, r *http.Request) {
	locationID, err := httputil.ParseUUIDParam(chi.URLParam(r, "location_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	locations, err := h.queries.GetLocationByID(r.Context(), locationID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "location not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, locations)
}

func (h *Handler) GetLocationsByPerson(w http.ResponseWriter, r *http.Request) {
	personID, err := httputil.ParseUUIDParam(chi.URLParam(r, "person_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	locations, err := h.queries.GetPersonLocations(r.Context(), personID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "location not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, locations)
}
