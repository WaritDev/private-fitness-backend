package usecases

import (
	"context"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserUseCase struct {
	repo repositories.UserRepo
}

func ProvideUserUseCase(repo repositories.UserRepo) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) ListUsers(ctx context.Context) ([]dbmodel.User, error) {
	return u.repo.List(ctx)
}
