package requests

import (
	
)

type CreateUserReq struct {
	Username        string    `json:"username" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	ImageURL    string    `json:"imageUrl" validate:"required,url"`
}