package rest

import "github.com/WaritDev/private-fitness-backend/domain/usecases"

type UserHandler struct {
    UC *usecases.UserUseCase
}

func ProvideUserHandler(uc *usecases.UserUseCase) *UserHandler {
    return &UserHandler{UC: uc}
}