package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"golangProject/internal/httpserver"
	"golangProject/internal/lib/logger/sl"
	"golangProject/internal/storage/postgresql"
	"log/slog"
	"os"
)

func main() {
	log := setupLogger()

	log.Info("Starting transaction-system application")
	log.Debug("debug messages are enabled")

	if err := godotenv.Load(); err != nil {
		log.Error("Error loading .env file ", sl.Error(err))
		os.Exit(1)
	}

	db, err := sql.Open("postgres", postgresql.ConnectionString())

	if err != nil {
		log.Error("Error connection in db ", sl.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	// TODO init NAT
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Error("Error connection NAT ", sl.Error(err))
		os.Exit(1)
	}
	defer nc.Close()

	database := postgresql.NewDatabase(db)

	server := httpserver.NewHTTPServer(database, nc)

	server.Start("8080")
}

func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
