package auth

import (
	"context"
	//"errors"

	"github.com/Nerzal/gocloak/v12"
	//"github.com/labstack/gommon/email"
	"fmt"

	"bytes"
	"encoding/json"
	
	"net/http"
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
// This method obtains an access token using LoginClient and then creates a user.
func (kc *KeycloakClient) RegisterUser(username, email, firstname, lastname, password string) error {
	ctx := context.Background()
	fmt.Printf(kc.ClientID)

	adminToken, err := kc.Client.LoginClient(ctx, kc.ClientID, kc.ClientSecret, kc.Realm)
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
	token := adminToken.AccessToken
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "http://localhost:8080/admin/realms/realm1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	_, err1 := client.Do(req)
	if err1 != nil {
		return err1
	}
	return err
}

// GetUsers retrieves all users from Keycloak
func (kc *KeycloakClient) GetUsers() ([]*gocloak.User, error) {
	ctx := context.Background()

	// Get admin token using client credentials
	token, err := kc.Client.LoginClient(
		ctx,
		kc.ClientID,
		kc.ClientSecret,
		kc.Realm,
	)
	//token, err := kc.Client.LoginAdmin(ctx, kc.AdminUsername, kc.AdminPassword, kc.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}

	// Get users with pagination parameters
	params := gocloak.GetUsersParams{
		First: gocloak.IntP(0),   // Offset
		Max:   gocloak.IntP(100), // Limit
	}

	users, err := kc.Client.GetUsers(ctx, token.AccessToken, kc.Realm, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	return users, nil
}

// GetUserByID retrieves a specific user by ID
func (kc *KeycloakClient) GetUserByID(userID string) (*gocloak.User, error) {
	ctx := context.Background()

	token, err := kc.Client.LoginClient(
		ctx,
		kc.ClientID,
		kc.ClientSecret,
		kc.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}

	user, err := kc.Client.GetUserByID(ctx, token.AccessToken, kc.Realm, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}
