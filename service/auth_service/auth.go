package auth_service

import (
	"go_gin_base/models"
)

type Auth struct {
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required,max=50"`
}

func (a *Auth) Check() (bool, error) {
	return models.CheckAuth(a.Username, a.Password)
}
