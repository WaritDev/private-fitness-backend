package psql

import (
    "context"

    "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserRepository struct {
    q *dbmodel.Queries
}

func ProvideUserRepository(q *dbmodel.Queries) *UserRepository {
    return &UserRepository{q: q}
}

func (r *UserRepository) List(ctx context.Context) ([]dbmodel.User, error) {
    return r.q.ListUsers(ctx)
}