package handlers


import (
    "net/http"

    "github.com/Twitter-Services-App/user-service/internal/core"
    "github.com/Twitter-Services-App/user-service/internal/models"
    "github.com/gin-gonic/gin"
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