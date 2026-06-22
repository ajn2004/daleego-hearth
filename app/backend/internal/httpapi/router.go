// internal/httpapi/router.go
package httpapi

import (
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/handlers/admin"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/handlers/mobile"
	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	"github.com/ajn2004/daleego-hearth/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(dbPool *pgxpool.Pool, AdminAPIKey string, PairingCodeSecret string) http.Handler {
	s := NewServer(dbPool)

	r := chi.NewRouter()

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Admin-API-Key"},
		AllowCredentials: true,
		MaxAge:           3000,
	}))

	adminHandler := admin.NewHandler(s.Queries, PairingCodeSecret)

	r.Get("/ping", s.ping)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AdminAPIKey(AdminAPIKey))
		r.Mount("/admin", adminHandler.Routes())
	})

	mobileHandler := mobile.NewHandler(s.Queries, s.DB, PairingCodeSecret)

	r.Group(func(r chi.Router) {
		// r.Use(middleware.DeviceAuth(s.Queries))
		r.Mount("/mobile", mobileHandler.Routes(s.Queries))
	})

	return r
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "pong",
	})
}
