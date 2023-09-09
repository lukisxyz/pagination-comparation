package main

import (
	"context"

	"github.com/flukis/pagination-comparation/db/storer"
	"github.com/flukis/pagination-comparation/internal/product"
	"github.com/flukis/pagination-comparation/utils/config"
	"github.com/flukis/pagination-comparation/utils/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig("")

	// database
	dbString := cfg.DBConfig.ConnString()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	// db storer
	productStorer := storer.NewProductStorer(pool)

	// services
	productSvc := product.NewService(productStorer)

	// router
	productRouter := product.NewRouter(productSvc)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// tight to chi router
	r.Mount("/api/product", productRouter.Routes())

	// Run server instance.
	log.Info().Msg("starting up server...")
	if err := server.Run(&cfg, r); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server Stopped")

}
