package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"natmap/internal/models"
)

// Cache 缓存接口
type Cache interface {
	GetMapping(tenantID, appID uint) (*models.MappingResponse, bool)
	SetMapping(tenantID, appID uint, mapping *models.MappingResponse)
	DeleteMapping(tenantID, appID uint)
	DeleteMappingByTenant(tenantID uint)
	DeleteMappingByApp(appID uint)
	Clear()
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	client *cache.Cache
	ttl    time.Duration
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache(defaultTTL, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		client: cache.New(defaultTTL, cleanupInterval),
		ttl:    defaultTTL,
	}
}

// mappingCacheKey 生成映射缓存键
func mappingCacheKey(tenantID, appID uint) string {
	return fmt.Sprintf("mapping:%d:%d", tenantID, appID)
}

// GetMapping 获取映射缓存
func (c *MemoryCache) GetMapping(tenantID, appID uint) (*models.MappingResponse, bool) {
	key := mappingCacheKey(tenantID, appID)
	value, found := c.client.Get(key)
	if !found {
		return nil, false
	}

	mapping, ok := value.(*models.MappingResponse)
	if !ok {
		return nil, false
	}

	return mapping, true
}

// SetMapping 设置映射缓存
func (c *MemoryCache) SetMapping(tenantID, appID uint, mapping *models.MappingResponse) {
	key := mappingCacheKey(tenantID, appID)
	c.client.Set(key, mapping, c.ttl)
}

// DeleteMapping 删除映射缓存
func (c *MemoryCache) DeleteMapping(tenantID, appID uint) {
	key := mappingCacheKey(tenantID, appID)
	c.client.Delete(key)
}

// DeleteMappingByTenant 删除租户下所有映射缓存
func (c *MemoryCache) DeleteMappingByTenant(tenantID uint) {
	// 遍历所有缓存项并删除匹配的
	items := c.client.Items()
	prefix := fmt.Sprintf("mapping:%d:", tenantID)
	for key := range items {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			c.client.Delete(key)
		}
	}
}

// DeleteMappingByApp 删除应用的所有映射缓存
func (c *MemoryCache) DeleteMappingByApp(appID uint) {
	// 由于缓存键是 tenantID:appID，需要遍历查找
	items := c.client.Items()
	suffix := fmt.Sprintf(":%d", appID)
	for key := range items {
		if len(key) > len(suffix) && key[len(key)-len(suffix):] == suffix {
			c.client.Delete(key)
		}
	}
}

// Clear 清空所有缓存
func (c *MemoryCache) Clear() {
	c.client.Flush()
}

// Stats 返回缓存统计信息
func (c *MemoryCache) Stats() map[string]interface{} {
	items := c.client.Items()
	return map[string]interface{}{
		"item_count": len(items),
		"default_ttl_seconds": c.ttl.Seconds(),
	}
}

// NoOpCache 空缓存实现（用于禁用缓存时）
type NoOpCache struct{}

func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

func (c *NoOpCache) GetMapping(tenantID, appID uint) (*models.MappingResponse, bool) {
	return nil, false
}

func (c *NoOpCache) SetMapping(tenantID, appID uint, mapping *models.MappingResponse) {}

func (c *NoOpCache) DeleteMapping(tenantID, appID uint) {}

func (c *NoOpCache) DeleteMappingByTenant(tenantID uint) {}

func (c *NoOpCache) DeleteMappingByApp(appID uint) {}

func (c *NoOpCache) Clear() {}

// SerializeMapping 序列化映射（用于调试）
func SerializeMapping(m *models.MappingResponse) string {
	data, _ := json.Marshal(m)
	return string(data)
}
