package database

import (
	"context"
	"log"

	"github.com/JSLee0920/golang-sveltekit-ecommerce/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(cfg *config.Config) *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Connected to DB")
	return pool
}
