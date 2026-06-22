package middleware

import (
	"log"
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/authkeys"
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/requestctx"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
)

const DeviceAPIKeyHeader = "X-Device-Api-Key"

func DeviceAuth(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deviceKey := r.Header.Get(DeviceAPIKeyHeader)
			if deviceKey == "" {
				log.Printf("Device tried to use API without proper key")
				response.WriteError(w, http.StatusUnauthorized, "missing device api key")
				return
			}

			prefix, err := authkeys.ExtractDeviceAPIKeyPrefix(deviceKey)
			if err != nil {
				log.Printf("device middleware used error: %v", err)
				response.WriteError(w, http.StatusUnauthorized, "could not extract prefix")
				return
			}

			deviceAPIKey, err := queries.GetActiveDeviceAPIKeysByPrefix(r.Context(), prefix)
			if err != nil {
				log.Printf("device middleware used error: %v", err)
				response.WriteError(w, http.StatusUnauthorized, "no device found")
				return
			}

			// mark device api key as used
			queries.MarkDeviceAPIKeyUsed(r.Context(), deviceAPIKey.ID)
			// mark device as last seen
			queries.UpdateDeviceLastSeen(r.Context(), deviceAPIKey.DeviceID)
			device, err := queries.GetDeviceByID(r.Context(), deviceAPIKey.DeviceID)
			if err != nil {
				log.Printf("device middleware used error: %v", err)
				response.WriteError(w, http.StatusNotFound, "device not found")
				return
			}
			ctx := requestctx.WithCurrentDevice(r.Context(), device)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
