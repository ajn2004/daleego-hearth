package httpapi

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

func NewServer(dbPool *pgxpool.Pool) *Server {
	return &Server{
		DB: dbPool,
	}
}
