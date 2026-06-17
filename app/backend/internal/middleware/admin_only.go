package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
)

const AdminAPIKeyHeader = "X-Admin-Api-Key"

func AdminAPIKey(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expectedKey == "" {
				response.WriteError(w, http.StatusInternalServerError, "admin api key is not configured")
				return
			}

			providedKey := r.Header.Get(AdminAPIKeyHeader)
			if providedKey == "" {
				response.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if subtle.ConstantTimeCompare([]byte(providedKey), []byte(expectedKey)) != 1 {
				response.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
