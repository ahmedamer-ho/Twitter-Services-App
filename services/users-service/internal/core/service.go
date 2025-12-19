package core

import (
	"context"

	"github.com/Twitter-Services-App/user-service/internal/models"
)

type AuthService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(username, password string) (string, error)
	GetUsers() ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetClientID() string
    GetClientSecret() string
    GetRealm() string
    Logout(ctx context.Context,token string) error
}
