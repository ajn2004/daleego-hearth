package httpapi

import (
	"github.com/ajn2004/daleego-hearth/backend/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB      *pgxpool.Pool
	Queries *db.Queries
}

func NewServer(dbPool *pgxpool.Pool) *Server {
	return &Server{
		DB:      dbPool,
		Queries: db.New(dbPool),
	}
}
