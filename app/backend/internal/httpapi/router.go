package httpapi

import (
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/httpapi/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(dbPool *pgxpool.Pool) http.Handler {
	s := NewServer(dbPool)

	r := chi.NewRouter()

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3000,
	}))

	r.Get("/ping", s.ping)
	return r
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "pong",
	})
}
