package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"natmap/internal/cache"
	"natmap/internal/models"
)

// AdminHandler 管理后台处理器
type AdminHandler struct {
	DB     *DB
	Cache  cache.Cache
	Logger *zap.Logger
}

// NewAdminHandler 创建管理后台处理器
func NewAdminHandler(db *DB, cache cache.Cache, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		DB:     db,
		Cache:  cache,
		Logger: logger,
	}
}

// ==================== 租户管理 ====================

// ListTenants 获取租户列表
// GET /api/admin?type=tenant
func (h *AdminHandler) ListTenants(c *gin.Context) {
	var tenants []models.Tenant
	result := h.DB.Conn.Order("created_at DESC").Find(&tenants)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var responses []models.AdminTenantResponse
	for _, t := range tenants {
		responses = append(responses, models.AdminTenantResponse{
			ID:         t.ID,
			TenantName: t.TenantName,
			CreatedAt:  t.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, responses)
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	TenantName string `json:"tenant_name" binding:"required"`
}

// CreateTenant 创建租户
// POST /api/admin?type=tenant
func (h *AdminHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant := models.Tenant{
		TenantName: req.TenantName,
	}

	if err := h.DB.Conn.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant created successfully",
		"data":    tenant,
	})
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	TenantName string `json:"tenant_name" binding:"required"`
}

// UpdateTenant 更新租户
// PUT /api/admin?type=tenant&id=xxx
func (h *AdminHandler) UpdateTenant(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tenant models.Tenant
	if err := h.DB.Conn.First(&tenant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	tenant.TenantName = req.TenantName
	if err := h.DB.Conn.Save(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant updated successfully",
	})
}

// DeleteTenant 删除租户
// DELETE /api/admin?type=tenant&id=xxx
func (h *AdminHandler) DeleteTenant(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	// 清除相关缓存
	if h.Cache != nil {
		tenantID, _ := strconv.ParseUint(id, 10, 32)
		h.Cache.DeleteMappingByTenant(uint(tenantID))
	}

	if err := h.DB.Conn.Delete(&models.Tenant{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant deleted successfully",
	})
}

// ==================== 应用管理 ====================

// ListApps 获取应用列表
// GET /api/admin?type=app
func (h *AdminHandler) ListApps(c *gin.Context) {
	var apps []models.App
	result := h.DB.Conn.Preload("Tenant").Order("created_at DESC").Find(&apps)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var responses []models.AdminAppResponse
	for _, a := range apps {
		tenantName := ""
		if a.Tenant != nil {
			tenantName = a.Tenant.TenantName
		}
		responses = append(responses, models.AdminAppResponse{
			ID:          a.ID,
			TenantID:    a.TenantID,
			TenantName:  tenantName,
			AppName:     a.AppName,
			Description: a.Description,
			CreatedAt:   a.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, responses)
}

// CreateAppRequest 创建应用请求
type CreateAppRequest struct {
	TenantID    uint   `json:"tenant_id" binding:"required"`
	AppName     string `json:"app_name" binding:"required"`
	Description string `json:"description"`
}

// CreateApp 创建应用
// POST /api/admin?type=app
func (h *AdminHandler) CreateApp(c *gin.Context) {
	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := models.App{
		TenantID:    req.TenantID,
		AppName:     req.AppName,
		Description: req.Description,
	}

	if err := h.DB.Conn.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "App created successfully",
		"data":    app,
	})
}

// UpdateAppRequest 更新应用请求
type UpdateAppRequest struct {
	AppName     string `json:"app_name" binding:"required"`
	Description string `json:"description"`
}

// UpdateApp 更新应用
// PUT /api/admin?type=app&id=xxx
func (h *AdminHandler) UpdateApp(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	var req UpdateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var app models.App
	if err := h.DB.Conn.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}

	app.AppName = req.AppName
	app.Description = req.Description
	if err := h.DB.Conn.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "App updated successfully",
	})
}

// DeleteApp 删除应用
// DELETE /api/admin?type=app&id=xxx
func (h *AdminHandler) DeleteApp(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	// 清除相关缓存
	if h.Cache != nil {
		appID, _ := strconv.ParseUint(id, 10, 32)
		h.Cache.DeleteMappingByApp(uint(appID))
	}

	if err := h.DB.Conn.Delete(&models.App{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "App deleted successfully",
	})
}

// ==================== 映射管理 ====================

// ListMappings 获取映射列表
// GET /api/admin?type=mapping
func (h *AdminHandler) ListMappings(c *gin.Context) {
	var mappings []models.Mapping
	result := h.DB.Conn.Preload("Tenant").Preload("App").
		Order("updated_at DESC").
		Limit(100).
		Find(&mappings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var responses []models.AdminMappingResponse
	for _, m := range mappings {
		tenantName := ""
		appName := ""
		if m.Tenant != nil {
			tenantName = m.Tenant.TenantName
		}
		if m.App != nil {
			appName = m.App.AppName
		}
		responses = append(responses, models.AdminMappingResponse{
			ID:         m.ID,
			TenantID:   m.TenantID,
			TenantName: tenantName,
			AppID:      m.AppID,
			AppName:    appName,
			PublicIP:   m.PublicIP,
			PublicPort: m.PublicPort,
			LocalIP:    m.LocalIP,
			LocalPort:  m.LocalPort,
			Protocol:   m.Protocol,
			UpdatedAt:  m.UpdatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, responses)
}

// CreateMappingRequest 创建映射请求
type CreateMappingRequest struct {
	TenantID   uint   `json:"tenant_id" binding:"required"`
	AppID      uint   `json:"app_id" binding:"required"`
	PublicIP   string `json:"public_ip" binding:"required,ip"`
	PublicPort uint   `json:"public_port" binding:"required,min=1,max=65535"`
	LocalIP    string `json:"local_ip" binding:"omitempty,ip"`
	LocalPort  uint   `json:"local_port" binding:"omitempty,min=1,max=65535"`
	Protocol   string `json:"protocol" binding:"omitempty,oneof=tcp udp"`
}

// CreateMapping 创建映射
// POST /api/admin?type=mapping
func (h *AdminHandler) CreateMapping(c *gin.Context) {
	var req CreateMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	protocol := req.Protocol
	if protocol == "" {
		protocol = "tcp"
	}

	mapping := models.Mapping{
		TenantID:   req.TenantID,
		AppID:      req.AppID,
		PublicIP:   req.PublicIP,
		PublicPort: req.PublicPort,
		LocalIP:    req.LocalIP,
		LocalPort:  req.LocalPort,
		Protocol:   protocol,
	}

	if err := h.DB.Conn.Create(&mapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mapping created successfully",
		"data":    mapping,
	})
}

// UpdateMappingRequest 更新映射请求
type UpdateMappingRequest struct {
	PublicIP   string `json:"public_ip" binding:"required,ip"`
	PublicPort uint   `json:"public_port" binding:"required,min=1,max=65535"`
	LocalIP    string `json:"local_ip" binding:"omitempty,ip"`
	LocalPort  uint   `json:"local_port" binding:"omitempty,min=1,max=65535"`
	Protocol   string `json:"protocol" binding:"omitempty,oneof=tcp udp"`
}

// UpdateMapping 更新映射
// PUT /api/admin?type=mapping&id=xxx
func (h *AdminHandler) UpdateMapping(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	var req UpdateMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mapping models.Mapping
	if err := h.DB.Conn.First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mapping not found"})
		return
	}

	// 清除旧缓存
	if h.Cache != nil {
		h.Cache.DeleteMapping(mapping.TenantID, mapping.AppID)
	}

	mapping.PublicIP = req.PublicIP
	mapping.PublicPort = req.PublicPort
	mapping.LocalIP = req.LocalIP
	mapping.LocalPort = req.LocalPort
	if req.Protocol != "" {
		mapping.Protocol = req.Protocol
	}

	if err := h.DB.Conn.Save(&mapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mapping updated successfully",
	})
}

// DeleteMapping 删除映射
// DELETE /api/admin?type=mapping&id=xxx
func (h *AdminHandler) DeleteMapping(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	// 先查询映射信息，用于清除缓存
	var mapping models.Mapping
	if err := h.DB.Conn.First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mapping not found"})
		return
	}

	// 清除缓存
	if h.Cache != nil {
		h.Cache.DeleteMapping(mapping.TenantID, mapping.AppID)
	}

	if err := h.DB.Conn.Delete(&mapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mapping deleted successfully",
	})
}
