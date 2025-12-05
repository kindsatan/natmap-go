package handlers

import (
    "net/http"

    "natmap-go/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AdminHandler struct { db *gorm.DB }
func NewAdminHandler(db *gorm.DB) *AdminHandler { return &AdminHandler{db: db} }

func (h *AdminHandler) ListUsers(c *gin.Context) {
    var users []models.User
    h.db.Select("id", "username", "email", "role", "is_active", "created_at", "updated_at", "last_login_at").Find(&users)
    c.JSON(http.StatusOK, users)
}

type setRoleReq struct { Role string `json:"role"` }
func (h *AdminHandler) SetUserRole(c *gin.Context) {
    var req setRoleReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Role == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"}); return }
    id := c.Param("id")
    if err := h.db.Model(&models.User{}).Where("id = ?", id).Update("role", req.Role).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"}); return }
    c.Status(http.StatusNoContent)
}

func (h *AdminHandler) ListPermissions(c *gin.Context) {
    var ps []models.Permission
    h.db.Find(&ps)
    c.JSON(http.StatusOK, ps)
}

type permReq struct {
    Role     string `json:"role"`
    Resource string `json:"resource"`
    Action   string `json:"action"`
    Allowed  *bool  `json:"allowed"`
}

func (h *AdminHandler) CreatePermission(c *gin.Context) {
    var req permReq
    if err := c.ShouldBindJSON(&req); err != nil || req.Role == "" || req.Resource == "" || req.Action == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    allowed := true
    if req.Allowed != nil { allowed = *req.Allowed }
    p := models.Permission{Role: req.Role, Resource: req.Resource, Action: req.Action, Allowed: allowed}
    if err := h.db.Create(&p).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"}); return }
    c.JSON(http.StatusCreated, p)
}

func (h *AdminHandler) DeletePermission(c *gin.Context) {
    id := c.Param("id")
    if err := h.db.Delete(&models.Permission{}, id).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"}); return }
    c.Status(http.StatusNoContent)
}
