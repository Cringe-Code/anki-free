package main

import (
	"anki/src/server"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getEnvValue(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while open .env file")
	}
	return os.Getenv(key)
}

func main() {

	logger := slog.Default()

	pgURL := getEnvValue("POSTGRES_CONN")

	if pgURL == "" {
		logger.Error("missed POSTGRES_CONN env")
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", pgURL)

	if err != nil {
		log.Fatalln(err)
	}

	serverAddress := getEnvValue("SERVER_ADDRESS")

	if serverAddress == "" {
		logger.Error("missed SERVER_ADDRESS")
		os.Exit(1)
	}

	signingKey := getEnvValue("SIGNING_KEY")

	if signingKey == "" {
		logger.Error("missed SIGNING_KEY")
		os.Exit(1)
	}

	s := server.NewServer(serverAddress, logger, db, signingKey)

	err = s.Start()

	if err != nil {
		logger.Error("server has been stopped", "error", err)
	}
}
