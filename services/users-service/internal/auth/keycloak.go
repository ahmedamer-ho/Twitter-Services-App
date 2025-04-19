package auth

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v12"
)

// KeycloakClient wraps the gocloak client and configuration for our realm.
type KeycloakClient struct {
	// Store the client as a pointer so it matches what NewClient returns.
	Client        *gocloak.GoCloak
	Realm         string
	ClientID      string
	ClientSecret  string
	AdminUsername string
	AdminPassword string
}

// NewKeycloakClient creates a new Keycloak client instance.
func NewKeycloakClient(url, realm, clientID, clientSecret, adminUsername, adminPassword string) *KeycloakClient {
	return &KeycloakClient{
		Client:        gocloak.NewClient(url), // NewClient returns *gocloak.GoCloak
		Realm:         realm,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
	}
}

// AuthenticateUser logs in a user and returns an access token.
func (kc *KeycloakClient) AuthenticateUser(username, password string) (string, error) {
	ctx := context.Background()
	token, err := kc.Client.Login(ctx, kc.ClientID, kc.ClientSecret, kc.Realm, username, password)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// RegisterUser creates a new user in Keycloak.
// This method obtains an admin token using LoginAdmin and then creates a user.
func (kc *KeycloakClient) RegisterUser(username,firstname,lastname, password, email string) error {
	ctx := context.Background()
	// Obtain an admin token.
	adminToken, err := kc.Client.LoginAdmin(ctx, kc.AdminUsername, kc.AdminPassword, kc.Realm)
	//adminToken, err := kc.Client.LoginAdmin(ctx, kc.AdminUsername, kc.AdminPassword, "master")
	if err != nil {
		return err
	}

	// Check if the access token is empty.
	if adminToken.AccessToken == "" {
		return errors.New("failed to obtain admin access token")
	}

	// Create a new user using CredentialRepresentation.
	user := gocloak.User{
		Username:  &username,
		Email:     &email,
		FirstName: &firstname,
		LastName:  &lastname,
		Enabled:   gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     &password,
				Temporary: gocloak.BoolP(false),
			},
		},
	}

	// CreateUser expects the admin token as a string.
	_, err = kc.Client.CreateUser(ctx, kc.Realm, adminToken.AccessToken, user)
	return err
}
