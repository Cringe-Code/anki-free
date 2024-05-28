package main

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	address string
	logger  *slog.Logger
	db      *sqlx.DB
}

func NewServer(address string, logger *slog.Logger, db *sqlx.DB) *Server {
	return &Server{
		address: address,
		logger:  logger,
		db:      db,
	}
}

func (s *Server) Start() error {
	return nil
}
