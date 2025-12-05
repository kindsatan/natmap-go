package handlers

import (
    "net/http"
    "time"

    "natmap-go/internal/auth"
    "natmap-go/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuthHandler struct {
    db          *gorm.DB
    jwtSecret   string
    tokenTTL    time.Duration
    refreshTTL  time.Duration
    bcryptCost  int
}

func NewAuthHandler(db *gorm.DB, secret string, ttl time.Duration, refreshTTL time.Duration, cost int) *AuthHandler {
    return &AuthHandler{db: db, jwtSecret: secret, tokenTTL: ttl, refreshTTL: refreshTTL, bcryptCost: cost}
}

type loginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req loginReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    var u models.User
    if err := h.db.Where("username = ?", req.Username).First(&u).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
        return
    }
    if !u.IsActive {
        c.JSON(http.StatusLocked, gin.H{"error": "user_locked"})
        return
    }
    if !auth.CheckPassword(u.PasswordHash, req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
        return
    }
    token, exp, err := auth.GenerateToken(u.ID, u.Username, u.Role, h.jwtSecret, h.tokenTTL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    raw, hash, err := auth.GenerateRefreshToken()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    rt := models.RefreshToken{UserID: u.ID, TokenHash: hash, ExpiresAt: time.Now().Add(h.refreshTTL)}
    if err := h.db.Create(&rt).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    now := time.Now()
    h.db.Model(&u).Updates(map[string]interface{}{"last_login_at": &now, "updated_at": now})
    c.JSON(http.StatusOK, gin.H{"token": token, "token_type": "Bearer", "expires_in": int(exp.Sub(time.Now()).Seconds()), "refresh_token": raw})
}

type registerReq struct {
    Username string  `json:"username"`
    Password string  `json:"password"`
    Email    *string `json:"email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req registerReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    var count int64
    h.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
    if count > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "username_taken"})
        return
    }
    hash, err := auth.HashPassword(req.Password, h.bcryptCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    now := time.Now().UTC()
    u := models.User{Username: req.Username, PasswordHash: hash, Email: req.Email, Role: "user", IsActive: true, CreatedAt: now, UpdatedAt: now}
    if err := h.db.Create(&u).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"id": u.ID, "username": u.Username})
}

type refreshReq struct {
    RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
    var req refreshReq
    if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    hash := auth.HashRefreshToken(req.RefreshToken)
    var rt models.RefreshToken
    if err := h.db.Where("token_hash = ?", hash).First(&rt).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh"})
        return
    }
    if rt.Revoked || time.Now().After(rt.ExpiresAt) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh"})
        return
    }
    var u models.User
    if err := h.db.First(&u, rt.UserID).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh"})
        return
    }
    token, exp, err := auth.GenerateToken(u.ID, u.Username, u.Role, h.jwtSecret, h.tokenTTL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    rawNew, hashNew, err := auth.GenerateRefreshToken()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    h.db.Model(&rt).Updates(map[string]interface{}{"revoked": true})
    _ = h.db.Create(&models.RefreshToken{UserID: u.ID, TokenHash: hashNew, ExpiresAt: time.Now().Add(h.refreshTTL)}).Error
    c.JSON(http.StatusOK, gin.H{"token": token, "token_type": "Bearer", "expires_in": int(exp.Sub(time.Now()).Seconds()), "refresh_token": rawNew})
}

func (h *AuthHandler) Logout(c *gin.Context) {
    var req refreshReq
    if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    hash := auth.HashRefreshToken(req.RefreshToken)
    if err := h.db.Model(&models.RefreshToken{}).Where("token_hash = ?", hash).Update("revoked", true).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    c.Status(http.StatusNoContent)
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
    uidVal, ok := c.Get("user_id")
    if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }
    uid, ok := uidVal.(uint)
    if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }
    if err := h.db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = ?", uid, false).Update("revoked", true).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
        return
    }
    c.Status(http.StatusNoContent)
}
