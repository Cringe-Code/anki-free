package server

import (
	"anki"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	address    string
	logger     *slog.Logger
	db         *sqlx.DB
	signingKey string
}

func NewServer(address string, logger *slog.Logger, db *sqlx.DB, signingKey string) *Server {
	return &Server{
		address:    address,
		logger:     logger,
		db:         db,
		signingKey: signingKey,
	}
}

func (s *Server) generateToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &anki.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(), // токен валиден в течение 12 часов
			IssuedAt:  time.Now().Unix(),
		},
		Login: login,
	})
	ResToken, err := token.SignedString([]byte(s.signingKey))
	return ResToken, err
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/ping", s.handlePing).Methods("GET")
	router.HandleFunc("/register", s.handleRegister).Methods("POST")
	router.HandleFunc("/auth", s.handlerSignIn).Methods("POST")
	router.HandleFunc("/packs/new", s.handleCreatePack).Methods("POST")

	s.logger.Info("server has been started", "address", s.address)

	err := http.ListenAndServe(s.address, router)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
