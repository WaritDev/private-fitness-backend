package sql

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
    rows, err := r.q.ListUsers(ctx)
    if err != nil {
        return nil, err
    }
    var users []dbmodel.User
    for _, row := range rows {
        users = append(users, dbmodel.User{
            Username:  row.Username,
            Email:     row.Email,
            CreatedAt: row.CreatedAt,
        })
    }
    return users, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetUserByUsernameRow, error) {
    return r.q.GetUserByUsername(ctx, username)
}