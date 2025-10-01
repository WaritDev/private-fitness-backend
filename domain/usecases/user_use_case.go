package usecases

import (
    "context"

    "github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/psql"
    "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserUseCase struct {
    repo *psql.UserRepository
}

func ProvideUserUseCase(repo *psql.UserRepository) *UserUseCase {
    return &UserUseCase{repo: repo}
}

func (u *UserUseCase) ListUsers(ctx context.Context) ([]dbmodel.User, error) {
    return u.repo.List(ctx)
}