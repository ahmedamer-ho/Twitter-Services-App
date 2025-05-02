package models

type User struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    FirstName string `json:"firstname"`
    LastName  string `json:"lastname"`
    Password  string `json:"password"`
    Enabled   bool   `json:"enabled"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type TokenResponse struct {
    Token string `json:"token"`
}