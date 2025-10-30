package usecases

import (
	"context"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/psql"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"fmt"
)

type AuthUseCase struct {
	authRepo *psql.AuthRepository
	userRepo *psql.UserRepository
}

func ProvideAuthUseCase(authRepo *psql.AuthRepository, userRepo *psql.UserRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo, userRepo: userRepo}
}

func (u *AuthUseCase) Login(ctx context.Context, req *dbmodel.UsersCredentials) (*dbmodel.UsersLoginRes, error) {
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, password is invalid")
	}

	token, err := u.authRepo.SignUsersAccessToken(&dbmodel.UsersPassport{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}
	res := &dbmodel.UsersLoginRes{
		AccessToken: token,
	}
	return res, nil
}
