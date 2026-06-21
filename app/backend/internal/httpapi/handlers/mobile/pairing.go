// internal/httpapi/handlers/mobile/pairing.go
package mobile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ajn2004/daleego-hearth/backend/internal/authkeys"
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	httputil "github.com/ajn2004/daleego-hearth/backend/internal/httpapi/utils"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) registerPairingRoutes(r chi.Router) {
	r.Post("/pairing", h.PairDevice)
}

type PairingRequest struct {
	PairCode  string `json:"pair_code"`
	Platform  string `json:"platform"`
	Model     string `json:"model"`
	ModelType string `json:"model_type"`
}

func (h *Handler) PairDevice(w http.ResponseWriter, r *http.Request) {

	var req PairingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	pairCode := strings.TrimSpace(req.PairCode)
	if pairCode == "" {
		response.WriteError(w, http.StatusBadRequest, "pair_code is required")
		return
	}

	// hash pairing code
	codeHash, err := httputil.HashValue(pairCode)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not hash code")
		return
	}

	// get pairing entry by code hash if entry exists, it is valid
	pairEntry, err := h.queries.GetValidDevicePairingCodeByHash(r.Context(), codeHash)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "could not find pairing")
		return
	}

	// get user associated
	person, err := h.queries.GetPersonByID(r.Context(), pairEntry.PersonID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "could not find person")
		return
	}
	model := strings.TrimSpace(req.Model)
	modelType := strings.TrimSpace(req.ModelType)
	if model == "" {
		response.WriteError(w, http.StatusBadRequest, "model required")
		return
	}

	// build Name of the form "User Name's Model ModelType"
	deviceName := fmt.Sprintf("%s's %s %s", person.DisplayName, model, modelType)

	devicePlatform, valid := parseDevicePlatform(strings.ToLower(strings.TrimSpace(req.Platform)))
	if !valid {
		response.WriteError(w, http.StatusBadRequest, "improper platform")
		return
	}

	// create device API key
	deviceAPIKey, err := authkeys.GenerateDeviceAPIKey()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not generate api key")
		return
	}

	tx, err := h.dbPool.Begin(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not begin transaction")
		return
	}
	defer tx.Rollback(r.Context())

	qtx := h.queries.WithTx(tx)
	// create device listing
	device, err := qtx.CreateDevice(r.Context(), db.CreateDeviceParams{
		PersonID: person.ID,
		Name:     deviceName,
		Platform: devicePlatform,
	})
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not create device")
		return
	}

	// create device API key entry
	_, err := qtx.CreateDeviceAPIKey(r.Context(), db.CreateDeviceAPIKeyParams{
		DeviceID:  device.ID,
		KeyHash:   deviceAPIKey.Hash,
		KeyPrefix: deviceAPIKey.Prefix,
	})

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not create api key record")
		return
	}

	_, err = qtx.MarkDevicePairingCodeUsed(r.Context(), db.MarkDevicePairingCodeUsedParams{
		DeviceID:      device.ID,
		PairingCodeID: pairEntry.ID,
	})
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not mark pairing code used")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "could not commit transaction")
		return
	}
	// return API key and device object to device
	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"api_key": deviceAPIKey.Plaintext,
		"device":  device,
	})
}

func parseDevicePlatform(value string) (db.DevicePlatform, bool) {
	switch value {
	case "android":
		return db.DevicePlatformAndroid, true
	case "ios":
		return db.DevicePlatformIos, true
	case "desktop":
		return db.DevicePlatformDesktop, true
	case "server":
		return db.DevicePlatformServer, true
	case "other":
		return db.DevicePlatformOther, true
	default:
		return "", false
	}
}
