package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"natmap/internal/config"
	"natmap/internal/middleware"
	"natmap/internal/models"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	DB     *DB
	Config *config.Config
	Logger *zap.Logger
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *DB, cfg *config.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		DB:     db,
		Config: cfg,
		Logger: logger,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

// Login 用户登录
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user models.User
	result := h.DB.Conn.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		h.Config.JWT.Secret,
		h.Config.JWT.ExpireHours,
	)
	if err != nil {
		h.Logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// GetCurrentUser 获取当前用户信息
// GET /api/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
		"role":     role,
	})
}

// InitDefaultUser 初始化默认管理员用户
func (h *AuthHandler) InitDefaultUser() error {
	var count int64
	h.DB.Conn.Model(&models.User{}).Count(&count)

	// 如果没有用户，创建默认管理员
	if count == 0 {
		h.Logger.Info("Creating default admin user")

		// 生成密码哈希
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := models.User{
			Username:     "admin",
			PasswordHash: string(hash),
			Role:         "admin",
		}

		if err := h.DB.Conn.Create(&admin).Error; err != nil {
			return err
		}

		h.Logger.Info("Default admin user created", zap.String("username", "admin"))
	}

	return nil
}

// ValidateToken 验证 JWT token
func (h *AuthHandler) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
