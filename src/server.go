package main

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc("/register", s.handleRegister).Methods("POST")

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
		s.logger.Error("error while parse json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while parse json"}`))
		return
	}

	var exists bool
	q := "select exists(select 1 from users where login=$1)"
	err = s.db.QueryRow(q, user.Login).Scan(&exists)

	if err != nil {
		s.logger.Error("error while checking exists", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while checking exists"}`))
		return
	}

	if exists {
		s.logger.Error("user with same login exists", "error")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"reason": "error while checking exists"}`))
		return
	}

	q = "insert into users(login) values ($1)"
	_, err = s.db.Exec(q, user.Login)

	if err != nil {
		s.logger.Error("error while insert into users table", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while insert into users table"}`))
		return
	}

	ans := anki.User{
		Login: user.Login,
	}

	userJson, _ := json.Marshal(ans)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"profile": %s}`, userJson)))
}
