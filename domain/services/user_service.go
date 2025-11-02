package services

import (
	"context"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserService struct {
	repo repositories.UserRepo
}

func ProvideUserService(repo repositories.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) ListUsers(ctx context.Context) ([]dbmodel.User, error) {
	return u.repo.List(ctx)
}

// CheckPhoneNumberExists checks if phone number is already registered (Q3S.1)
func (u *UserService) CheckPhoneNumberExistsUser(ctx context.Context, phoneNumber string) (bool, error) {
	return u.repo.CheckPhoneNumberExistsUser(ctx, phoneNumber)
}

// CheckGmailExists checks if gmail is already registered (Q3S.2)
func (u *UserService) CheckGmailExistsUser(ctx context.Context, gmail string) (bool, error) {
	return u.repo.CheckGmailExistsUser(ctx, gmail)
}

func (u *UserService) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	return u.repo.CheckUsernameExists(ctx, username)
}


