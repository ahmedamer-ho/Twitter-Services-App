package handlers

import (
    "net/http"

    "github.com/Twitter-Services-App/user-service/internal/core"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    service core.AuthService
}

func NewUserHandler(service core.AuthService) *UserHandler {
    return &UserHandler{service: service}
}


func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.GET("/users", h.getUsers)
    router.GET("/users/:id", h.getUserByID)
}

func (h *UserHandler) getUsers(c *gin.Context) {
    users, err := h.service.GetUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) getUserByID(c *gin.Context) {
    userID := c.Param("id")
    user, err := h.service.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"user": user})
}