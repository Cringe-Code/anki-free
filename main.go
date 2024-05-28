package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
)

func main() {
	logger := slog.Default()
	pgURL := os.Getenv("POSTGRES_CONN")

	if pgURL == "" {
		logger.Error("missed POSTGRES_CONN env")
		os.Exit(1)
	}

	db, err := sqlx.Connect("pqx", pgURL)

	if err != nil {
		log.Fatalln(err)
	}

	_ = db

	serverAddress := os.Getenv("SERVER_ADDRESS")

	if serverAddress == "" {
		logger.Error("missed SERVER_ADDRES")
		os.Exit(1)
	}

	s := NewServer(serverAddress, logger, db)

	err = s.Start()

	if err != nil {
		logger.Error("server has been stopped", "error", err)
	}

	logger.Info("server has been started")

}
