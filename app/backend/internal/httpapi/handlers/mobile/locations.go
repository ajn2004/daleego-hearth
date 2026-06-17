// internal/httpapi/handlers/mobile/locations.go
package mobile

import "net/http"

func (h *Handler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	// device identity comes from middleware
	// parse location payload
	// call h.LocationService.CreateLocation(...)
}
