package models

import (
	"time"
	"gorm.io/gorm"
)

// Tenant 租户模型
type Tenant struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantName  string         `gorm:"column:tenant_name;type:varchar(255);not null;uniqueIndex:uk_tenant_name" json:"tenant_name"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Apps        []App          `gorm:"foreignKey:TenantID;references:ID" json:"apps,omitempty"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// App 应用模型
type App struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    uint           `gorm:"column:tenant_id;not null;index:idx_tenant_id;uniqueIndex:uk_tenant_app" json:"tenant_id"`
	AppName     string         `gorm:"column:app_name;type:varchar(255);not null;uniqueIndex:uk_tenant_app" json:"app_name"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID;references:ID" json:"tenant,omitempty"`
}

func (App) TableName() string {
	return "apps"
}

// Mapping 映射模型
type Mapping struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    uint           `gorm:"column:tenant_id;not null;index:idx_tenant_id;uniqueIndex:uk_tenant_app" json:"tenant_id"`
	AppID       uint           `gorm:"column:app_id;not null;index:idx_app_id;uniqueIndex:uk_tenant_app" json:"app_id"`
	PublicIP    string         `gorm:"column:public_ip;type:varchar(45);not null" json:"public_ip"`
	PublicPort  uint           `gorm:"column:public_port;not null" json:"public_port"`
	LocalIP     string         `gorm:"column:local_ip;type:varchar(45)" json:"local_ip"`
	LocalPort   uint           `gorm:"column:local_port" json:"local_port"`
	Protocol    string         `gorm:"column:protocol;type:varchar(10);default:'tcp'" json:"protocol"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID;references:ID" json:"tenant,omitempty"`
	App         *App           `gorm:"foreignKey:AppID;references:ID" json:"app,omitempty"`
}

func (Mapping) TableName() string {
	return "mappings"
}

// User 用户模型
type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username" json:"username"`
	PasswordHash string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"column:role;type:enum('admin','user');default:'user'" json:"role"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// MappingResponse 映射查询响应（兼容原 API 格式）
type MappingResponse struct {
	PublicIP   string `json:"public_ip"`
	PublicPort uint   `json:"public_port"`
	UpdatedAt  string `json:"updated_at"`
	Cache      string `json:"_cache,omitempty"`
}

// ToResponse 转换为 API 响应格式
func (m *Mapping) ToResponse() *MappingResponse {
	// 转换为北京时间 (UTC+8)
	beijingTime := m.UpdatedAt.Add(8 * time.Hour)
	return &MappingResponse{
		PublicIP:   m.PublicIP,
		PublicPort: m.PublicPort,
		UpdatedAt:  beijingTime.Format("2006-01-02 15:04:05"),
	}
}

// AdminMappingResponse 管理后台映射响应
type AdminMappingResponse struct {
	ID          uint      `json:"id"`
	TenantID    uint      `json:"tenant_id"`
	TenantName  string    `json:"tenant_name"`
	AppID       uint      `json:"app_id"`
	AppName     string    `json:"app_name"`
	PublicIP    string    `json:"public_ip"`
	PublicPort  uint      `json:"public_port"`
	LocalIP     string    `json:"local_ip"`
	LocalPort   uint      `json:"local_port"`
	Protocol    string    `json:"protocol"`
	UpdatedAt   string    `json:"updated_at"`
}

// AdminAppResponse 管理后台应用响应
type AdminAppResponse struct {
	ID          uint      `json:"id"`
	TenantID    uint      `json:"tenant_id"`
	TenantName  string    `json:"tenant_name"`
	AppName     string    `json:"app_name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
}

// AdminTenantResponse 管理后台租户响应
type AdminTenantResponse struct {
	ID         uint      `json:"id"`
	TenantName string    `json:"tenant_name"`
	CreatedAt  string    `json:"created_at"`
}
