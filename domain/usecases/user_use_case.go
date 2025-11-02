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
func (u *UserUseCase) CheckPhoneNumberExistsUser(ctx context.Context, phoneNumber string) (bool, error) {
	return u.repo.CheckPhoneNumberExistsUser(ctx, phoneNumber)
}

// CheckGmailExists checks if gmail is already registered (Q3S.2)
func (u *UserUseCase) CheckGmailExistsUser(ctx context.Context, gmail string) (bool, error) {
	return u.repo.CheckGmailExistsUser(ctx, gmail)
}

func (u *UserUseCase) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	return u.repo.CheckUsernameExists(ctx, username)
}


