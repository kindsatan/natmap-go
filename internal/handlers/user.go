package handlers

import (
    "net/http"

    "natmap-go/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type UserHandler struct {
    db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler { return &UserHandler{db: db} }

func (h *UserHandler) Me(c *gin.Context) {
    uidVal, ok := c.Get("user_id")
    if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }
    uid, ok := uidVal.(uint)
    if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }
    var u models.User
    if err := h.db.First(&u, uid).Error; err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }
    c.JSON(http.StatusOK, gin.H{"id": u.ID, "username": u.Username, "email": u.Email, "last_login_at": u.LastLoginAt})
}
