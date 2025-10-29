package psql

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type AuthRepository struct {
	q *dbmodel.Queries
}

func ProvideAuthRepository(q *dbmodel.Queries) *AuthRepository {
	return &AuthRepository{q: q}
}

func (r *AuthRepository) SignUsersAccessToken(req *dbmodel.UsersPassport) (string, error) {
	claims := dbmodel.UsersClaims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "access_token",
			Subject:   "users_access_token",
			ID:        uuid.NewString(),
			Audience:  []string{"users"},
		},
	}

	mySigningKey := os.Getenv("JWT_SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return ss, nil
}
