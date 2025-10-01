package db

import (
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

func ProvideQueries(pool *pgxpool.Pool) *dbmodel.Queries {
    return dbmodel.New(pool)
}