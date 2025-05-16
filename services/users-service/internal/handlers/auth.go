package handlers


import (
    "net/http"

    "github.com/Twitter-Services-App/user-service/internal/core"
    "github.com/Twitter-Services-App/user-service/internal/models"
    "github.com/gin-gonic/gin"
    "strings" 
)

type AuthHandler struct {
    service core.AuthService
}

func NewAuthHandler(service core.AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}


func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
    router.GET("/", h.home)
    
    authGroup := router.Group("/auth")
    {
        authGroup.POST("/register", h.register)
        authGroup.POST("/login", h.login)
        authGroup.POST("/logout", h.logout)
    }
}

func (h *AuthHandler) home(c *gin.Context) {
    c.JSON(200, gin.H{"message": "users-services"})
}

func (h *AuthHandler) register(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    if err := h.service.RegisterUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    token, err := h.service.AuthenticateUser(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
        return
    }

    c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}
func (h *AuthHandler) logout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
        return
    }

    // Extract Bearer token
    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
        return
    }

    token := tokenParts[1]
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
        return
    }

    // Perform logout
    if err := h.service.Logout(c,token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "logout failed",
            "details": err.Error(),
        })
        return
    }

    // Clear client-side data
    c.Header("Clear-Site-Data", `"cookies", "storage", "cache"`)
    c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}