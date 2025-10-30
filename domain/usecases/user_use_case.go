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

// CheckPhoneNumberExists checks if phone number is already registered (Q3S.1)
func (u *UserUseCase) CheckPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	return u.repo.CheckPhoneNumberExists(ctx, phoneNumber)
}

// CheckGmailExists checks if gmail is already registered (Q3S.2)
func (u *UserUseCase) CheckGmailExists(ctx context.Context, gmail string) (bool, error) {
	return u.repo.CheckGmailExists(ctx, gmail)
}
