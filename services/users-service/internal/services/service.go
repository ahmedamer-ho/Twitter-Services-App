package auth

import (

    "github.com/Twitter-Services-App/user-service/internal/core"
    "github.com/Twitter-Services-App/user-service/internal/models"
)

type KeycloakService struct {
    client *KeycloakClient
}

func NewKeycloakService(client *KeycloakClient) core.AuthService {
    return &KeycloakService{client: client}
}

func (s *KeycloakService) RegisterUser(user *models.User) error {
    return s.client.RegisterUser(user.Username, user.Email, user.FirstName, user.LastName, user.Password)
}

func (s *KeycloakService) AuthenticateUser(username, password string) (string, error) {
    return s.client.AuthenticateUser(username, password)
}

func (s *KeycloakService) GetUsers() ([]*models.User, error) {
    gUsers, err := s.client.GetUsers()
    if err != nil {
        return nil, err
    }

    users := make([]*models.User, len(gUsers))
    for i, gUser := range gUsers {
        users[i] = &models.User{
        ID:        safeString(gUser.ID),
        Username:  safeString(gUser.Username),
        Email:     safeString(gUser.Email),
        FirstName: safeString(gUser.FirstName),
        LastName:  safeString(gUser.LastName),
        Enabled:   safeBool(gUser.Enabled),
        }
    }
    return users, nil
}

func (s *KeycloakService) GetUserByID(id string) (*models.User, error) {
    gUser, err := s.client.GetUserByID(id)
    if err != nil {
        return nil, err
    }

    return &models.User{
        ID:        *gUser.ID,
        Username:  *gUser.Username,
        Email:     *gUser.Email,
        FirstName: *gUser.FirstName,
        LastName:  *gUser.LastName,
        Enabled:   *gUser.Enabled,
    }, nil
}

func safeString(ptr *string) string {
    if ptr != nil {
        return *ptr
    }
    return ""
}

func safeBool(ptr *bool) bool {
    if ptr != nil {
        return *ptr
    }
    return false
}

