package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marmiz/pgrest-explorations/internal/db"
)

type config struct {
	port int
	env  string
}

type application struct {
	config  config
	logger  *slog.Logger
	queries *db.Queries
}

func main() {
	ctx := context.Background()

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("unable to connect to database", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	app := &application{
		config:  cfg,
		logger:  logger,
		queries: queries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("/v1/user", app.getUser)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "port", cfg.port, "env", cfg.env)

	err = server.ListenAndServe()

	logger.Error("server shutdown", err)

}
