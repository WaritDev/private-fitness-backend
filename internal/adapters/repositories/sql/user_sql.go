package sql

import (
    "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserRepository struct {
    q *dbmodel.Queries
}

func ProvideUserRepository(q *dbmodel.Queries) *UserRepository {
    return &UserRepository{q: q}
}