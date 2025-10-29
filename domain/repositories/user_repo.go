package repositories

import (
	"context"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserRepo interface {
	ListUsers(ctx context.Context) ([]*dbmodel.User, error)
	GetUserByUsername(ctx context.Context) (dbmodel.GetUserByUsernameRow, error)
}