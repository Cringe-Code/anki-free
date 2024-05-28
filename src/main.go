package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	logger := slog.Default()

	pgURL := os.Getenv("POSTGRES_CONN")

	if pgURL == "" {
		err := godotenv.Load(".env")
		if err != nil {
			logger.Error("cant open .env file")
			os.Exit(1)
		}

		pgURL := os.Getenv("POSTGRES_CONN")

		if pgURL == "" {
			logger.Error("missed POSTGRES_CONN env")
			os.Exit(1)
		}
	}

	db, err := sqlx.Connect("postgres", pgURL)

	if err != nil {
		log.Fatalln(err)
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")

	if serverAddress == "" {
		err = godotenv.Load(".env")

		if err != nil {
			logger.Error("cant open .env file")
			os.Exit(1)
		}

		serverAddress = os.Getenv("SERVER_ADDRESS")

		if serverAddress == "" {
			logger.Error("missed SERVER_ADDRESS")
			os.Exit(1)
		}
	}

	s := NewServer(serverAddress, logger, db)

	err = s.Start()

	if err != nil {
		logger.Error("server has been stopped", "error", err)
	}
}
