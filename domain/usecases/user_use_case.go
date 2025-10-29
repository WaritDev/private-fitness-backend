package usecases

import (

    "github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/sql"
)

type UserUseCase struct {
    repo *sql.UserRepository
}

func ProvideUserUseCase(repo *sql.UserRepository) *UserUseCase {
    return &UserUseCase{repo: repo}
}
