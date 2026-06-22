// internal/httpapi/handlers/mobile/locations.go
package mobile

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/requestctx"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) registerLocationRoutes(r chi.Router) {
	r.Post("/locations", h.ReportLocation)
}

type LocationRequest struct {
	Latitude   float64   `json:"lat"`
	Longitude  float64   `json:"lng"`
	Accuracy   float64   `json:"accuracy"`
	RecordedAt time.Time `json:"recorded_at"`
}

func (h *Handler) ReportLocation(w http.ResponseWriter, r *http.Request) {
	// device identity comes from middleware
	// parse location payload
	// call h.LocationService.CreateLocation(...)
	var req LocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	device, ok := requestctx.CurrentDevice(r)
	if !ok {
		log.Printf("Failed to find device in context")
		response.WriteError(w, http.StatusInternalServerError, "could not find device in context")
		return
	}
	location, err := h.queries.CreateLocation(r.Context(), db.CreateLocationParams{
		Lat:        req.Latitude,
		Lng:        req.Longitude,
		Accuracy:   req.Accuracy,
		RecordedAt: pgtype.Timestamptz{Valid: true, Time: req.RecordedAt},
		DeviceID:   device.ID,
		PersonID:   device.PersonID,
	})
	if err != nil {
		log.Printf("mark pairing used error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "could not record location")
		return
	}

	response.WriteJSON(w, http.StatusCreated, location)
}
