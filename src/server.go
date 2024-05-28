package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"anki"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/ping", s.handlePing).Methods("GET")

	s.logger.Info("server has been started", "address", s.address)

	err := http.ListenAndServe(s.address, router)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user anki.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {

	}
}
