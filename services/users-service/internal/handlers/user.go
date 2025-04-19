package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/Twitter-Services-App/user-service/internal/db"
    "gorm.io/gorm"
)

type UserHandler struct {
    DB *gorm.DB
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
    var user db.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": "Invalid input"})
        return
    }

    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(201, user)
}
