package core

import "github.com/Twitter-Services-App/user-service/internal/models"

type AuthService interface {
	RegisterUser(user models.User) error
	AuthenticateUser(username, password string) (string, error)
	GetUsers() ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
}
