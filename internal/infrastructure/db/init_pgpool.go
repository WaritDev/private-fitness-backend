package db

import (
	"context"
	"fmt"
	"log"

	"github.com/WaritDev/private-fitness-backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProvidePgxPool(ctx context.Context, config *config.Config) *pgxpool.Pool {
	dbPort := config.PostgresPort
	if config.UseDocker {
		dbPort = "5432"
	}
	postgresURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config.PostgresHost, dbPort, config.PostgresUser, config.PostgresPassword, config.PostgresDB)

	connConfig, err := pgxpool.ParseConfig(postgresURI)
	if err != nil {
		panic(err)
	}

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		panic(err)
	}

	// Connection Teardown
	conn, err := pgxPool.Acquire(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Release()

	if err := conn.Conn().Ping(ctx); err != nil {
		panic(err)
	}

	log.Println("🫙 Connected to Postgres")

	return pgxPool
}

func GetPgPool() *pgxpool.Pool {
	ctx := context.Background()
	config := config.ProvideConfig()
	pool := ProvidePgxPool(ctx, config)
	return pool
}