package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marmiz/pgrest-explorations/internal/db"
)

type config struct {
	port int
	env  string
	jwt  struct {
		secret string
	}
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

	s, varErr := os.LookupEnv("JWT_SECRET")
	if !varErr {
		fmt.Println("JWT_SECRET not set")
		os.Exit(1)
	}

	cfg.jwt.secret = s

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Get("/v1/user", app.getUser)
	r.Post("/v1/auth/token", app.createAuthTokenHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "port", cfg.port, "env", cfg.env)

	err = server.ListenAndServe()

	logger.Error("server shutdown", err)

}
