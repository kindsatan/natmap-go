package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"natmap/internal/cache"
	"natmap/internal/models"
)

// MappingHandler 映射处理器
type MappingHandler struct {
	DB     *DB
	Cache  cache.Cache
	Logger *zap.Logger
}

// NewMappingHandler 创建映射处理器
func NewMappingHandler(db *DB, cache cache.Cache, logger *zap.Logger) *MappingHandler {
	return &MappingHandler{
		DB:     db,
		Cache:  cache,
		Logger: logger,
	}
}

// GetMappingRequest 获取映射请求参数
type GetMappingRequest struct {
	TenantID uint `form:"tenant_id" binding:"required"`
	AppID    uint `form:"app_id" binding:"required"`
}

// GetMapping 获取映射信息
// GET /api/get?tenant_id=xxx&app_id=xxx
func (h *MappingHandler) GetMapping(c *gin.Context) {
	var req GetMappingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing tenant_id or app_id parameter",
		})
		return
	}

	// 1. 先尝试从缓存获取
	if h.Cache != nil {
		if cached, found := h.Cache.GetMapping(req.TenantID, req.AppID); found {
			cached.Cache = "HIT"
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	// 2. 缓存未命中，查询数据库
	var mapping models.Mapping
	result := h.DB.Conn.Where("tenant_id = ? AND app_id = ?", req.TenantID, req.AppID).
		Order("updated_at DESC").
		First(&mapping)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}

	// 3. 转换为响应格式
	resp := mapping.ToResponse()

	// 4. 写入缓存
	if h.Cache != nil {
		h.Cache.SetMapping(req.TenantID, req.AppID, resp)
	}

	resp.Cache = "MISS"
	c.JSON(http.StatusOK, resp)
}

// UpdateMappingAPIRequest 更新映射请求 (API 版本)
type UpdateMappingAPIRequest struct {
	App       string `json:"app" binding:"required"`
	IP        string `json:"ip" binding:"required,ip"`
	Port      uint   `json:"port" binding:"required,min=1,max=65535"`
	Proto     string `json:"proto" binding:"omitempty,oneof=tcp udp"`
	LocalIP   string `json:"local_ip" binding:"omitempty,ip"`
	LocalPort uint   `json:"local_port" binding:"omitempty,min=1,max=65535"`
}

// UpdateMapping 更新映射信息
// POST /api/update
func (h *MappingHandler) UpdateMapping(c *gin.Context) {
	// 从上下文获取租户ID（由中间件设置）
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req UpdateMappingAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tenantIDUint := tenantID.(uint)

	// 查询应用
	var app models.App
	result := h.DB.Conn.Where("tenant_id = ? AND app_name = ?", tenantIDUint, req.App).First(&app)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "app not found",
		})
		return
	}

	// 使用事务保证原子性
	err := h.DB.Conn.Transaction(func(tx *gorm.DB) error {
		// 删除旧映射
		if err := tx.Where("tenant_id = ? AND app_id = ?", tenantIDUint, app.ID).Delete(&models.Mapping{}).Error; err != nil {
			return err
		}

		// 创建新映射
		protocol := req.Proto
		if protocol == "" {
			protocol = "tcp"
		}

		mapping := models.Mapping{
			TenantID:   tenantIDUint,
			AppID:      app.ID,
			PublicIP:   req.IP,
			PublicPort: req.Port,
			LocalIP:    req.LocalIP,
			LocalPort:  req.LocalPort,
			Protocol:   protocol,
		}

		if err := tx.Create(&mapping).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		h.Logger.Error("Failed to update mapping", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update mapping",
		})
		return
	}

	// 清除缓存
	if h.Cache != nil {
		h.Cache.DeleteMapping(tenantIDUint, app.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mapping updated successfully",
	})
}
