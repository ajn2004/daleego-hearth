// httpapi/handlers/admin/pairings.go
package admin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	httputil "github.com/ajn2004/daleego-hearth/backend/internal/httpapi/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

/*
API Structure

GET /admin/pairings

	Description:
	  Lists all pairing codes and their status.
	  Useful for admin/debug visibility.

POST /admin/pairings

	Payload:
	  {
	    "person_id": "uuid"
	  }
	Description:
	  Creates a new device pairing code for a person.
	Returns:
	  Pairing code response object, including the plaintext code.

GET /admin/pairings/{pairing_id}

	Description:
	  Gets one pairing code by id.
	Returns:
	  Pairing code database object/status.

POST /admin/pairings/{pairing_id}/revoke

	Description:
	  Revokes or expires a pairing code.
	  This should not delete a device.
	Returns:
	  Updated pairing code object.
*/
func (h *Handler) registerPairingRoutes(r chi.Router) {
	r.Get("/pairings", h.ListPairings)
	r.Post("/pairings", h.CreatePairing)
	r.Get("/pairings/expired", h.ListExpiredPairings)
	r.Get("/pairings/{pairing_id}", h.GetPairing)
	r.Post("/pairings/{pairing_id}/revoke", h.RevokePairing)
}

func (h *Handler) ListPairings(w http.ResponseWriter, r *http.Request) {
	pairings, err := h.queries.GetActivePairings(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not list active pairings")
		return
	}
	response.WriteJSON(w, http.StatusOK, pairings)
}

func (h *Handler) ListExpiredPairings(w http.ResponseWriter, r *http.Request) {
	pairings, err := h.queries.GetExpiredPairings(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not list active pairings")
		return
	}
	response.WriteJSON(w, http.StatusOK, pairings)
}

type PairingRequest struct {
	PersonID string `json:"person_id"`
}

func (h *Handler) CreatePairing(w http.ResponseWriter, r *http.Request) {
	var req PairingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "request needs person_id")
		return
	}
	personID, err := httputil.ParseUUIDParam(req.PersonID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "failed to make person_id uuid")
		return
	}

	genCode, err := generateRandomPairing()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create code")
		return
	}

	hashedCode, err := httputil.HashValue(genCode)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to hash code")
		return
	}

	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	pairing, err := h.queries.CreateDevicePairingCode(r.Context(), db.CreateDevicePairingCodeParams{
		PersonID:  personID,
		CodeHash:  hashedCode,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to make pairing")
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"pairing_code": genCode,
		"pairing":      pairing,
	})
}

func (h *Handler) GetPairing(w http.ResponseWriter, r *http.Request) {
	pairingID, err := httputil.ParseUUIDParam(chi.URLParam(r, "pairing_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "failed to get id")
		return
	}

	pairing, err := h.queries.GetPairingCodeByID(r.Context(), pairingID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "failed to find pairing code")
		return
	}
	response.WriteJSON(w, http.StatusOK, pairing)
}

func (h *Handler) RevokePairing(w http.ResponseWriter, r *http.Request) {
	pairingID, err := httputil.ParseUUIDParam(chi.URLParam(r, "pairing_id"))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get id")
		return
	}

	pairing, err := h.queries.RevokePairingCode(r.Context(), pairingID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "failed to find revoke code")
		return
	}
	response.WriteJSON(w, http.StatusOK, pairing)
}

func generateRandomPairing() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}

	code := int(n.Int64())
	return fmt.Sprintf("%03d-%03d", code/1000, code%1000), nil
}
