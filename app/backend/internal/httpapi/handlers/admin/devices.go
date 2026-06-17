// httpapi/handlers/admin/devices.go
package admin

import (
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	httputil "github.com/ajn2004/daleego-hearth/backend/internal/httpapi/utils"
	"github.com/go-chi/chi/v5"
)

/*
API Structure

GET /admin/devices

	Description:
	  Lists all devices.

GET /admin/devices/{device_id}

	Description:
	  Gets one device by ID.

DELETE /admin/devices/{device_id}

	Description:
	  Soft-deletes/deactivates a device.

POST /admin/devices/{device_id}/revoke-keys

	Description:
	  Revokes active API keys for a device.
*/
func (h *Handler) registerDeviceRoutes(r chi.Router) {
	r.Get("/devices", h.ListDevices)
	r.Get("/devices/{device_id}", h.GetDevice)
	r.Delete("/devices/{device_id}", h.DeleteDevice)
	r.Post("/devices/{device_id}/revoke-keys", h.RevokeDeviceKeys)
}

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.queries.GetAllDevices(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "unable to list devices")
		return
	}
	response.WriteJSON(w, http.StatusOK, devices)
}
func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, err := httputil.ParseUUIDParam(chi.URLParam(r, "device_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "unable to find device")
		return
	}

	device, err := h.queries.GetDeviceByID(r.Context(), deviceID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "unable to find device")
		return
	}
	response.WriteJSON(w, http.StatusOK, device)
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, err := httputil.ParseUUIDParam(chi.URLParam(r, "device_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "unable to find device")
		return
	}

	device, err := h.queries.DeleteDevice(r.Context(), deviceID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "unable to delete device")
		return
	}
	response.WriteJSON(w, http.StatusOK, device)
}
func (h *Handler) RevokeDeviceKeys(w http.ResponseWriter, r *http.Request) {
	deviceID, err := httputil.ParseUUIDParam(chi.URLParam(r, "device_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "unable to find device")
		return
	}

	keys, err := h.queries.RevokeDeviceAPIKeys(r.Context(), deviceID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "unable to revoke device keys")
		return
	}
	response.WriteJSON(w, http.StatusOK, keys)
}
