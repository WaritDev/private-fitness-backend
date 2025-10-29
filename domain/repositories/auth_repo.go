package repositories

import (
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type AuthRepo interface {
	SignUsersAccessToken(req *dbmodel.UsersPassport) (string, error)
}

