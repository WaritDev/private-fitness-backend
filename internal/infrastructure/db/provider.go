package db

import (
	"database/sql"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

func ProvideQueries(db *sql.DB) *dbmodel.Queries {
	return dbmodel.New(db)
}