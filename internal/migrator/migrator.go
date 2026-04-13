package migrator

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"natmap/internal/models"
)

// D1Record D1 数据库记录格式
type D1Record struct {
	ID          uint      `json:"id"`
	TenantID    uint      `json:"tenant_id"`
	TenantName  string    `json:"tenant_name"`
	AppID       uint      `json:"app_id"`
	AppName     string    `json:"app_name"`
	Description string    `json:"description"`
	PublicIP    string    `json:"public_ip"`
	PublicPort  uint      `json:"public_port"`
	LocalIP     string    `json:"local_ip"`
	LocalPort   uint      `json:"local_port"`
	Protocol    string    `json:"protocol"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Migrator 数据迁移器
type Migrator struct {
	DB         *gorm.DB
	D1Endpoint string
	D1Token    string
}

// NewMigrator 创建迁移器
func NewMigrator(db *gorm.DB, d1Endpoint, d1Token string) *Migrator {
	return &Migrator{
		DB:         db,
		D1Endpoint: d1Endpoint,
		D1Token:    d1Token,
	}
}

// MigrateFromD1 从 D1 迁移数据
func (m *Migrator) MigrateFromD1() error {
	fmt.Println("Starting migration from Cloudflare D1 to MySQL...")

	// 迁移租户
	if err := m.migrateTenants(); err != nil {
		return fmt.Errorf("failed to migrate tenants: %w", err)
	}

	// 迁移应用
	if err := m.migrateApps(); err != nil {
		return fmt.Errorf("failed to migrate apps: %w", err)
	}

	// 迁移映射
	if err := m.migrateMappings(); err != nil {
		return fmt.Errorf("failed to migrate mappings: %w", err)
	}

	fmt.Println("Migration completed successfully!")
	return nil
}

// migrateTenants 迁移租户数据
func (m *Migrator) migrateTenants() error {
	fmt.Println("Migrating tenants...")

	// 从 D1 获取租户数据
	tenants, err := m.fetchD1Tenants()
	if err != nil {
		return err
	}

	for _, t := range tenants {
		tenant := models.Tenant{
			ID:         t.ID,
			TenantName: t.TenantName,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}

		// 使用 INSERT IGNORE 避免重复
		if err := m.DB.FirstOrCreate(&tenant, models.Tenant{ID: t.ID}).Error; err != nil {
			return fmt.Errorf("failed to create tenant %d: %w", t.ID, err)
		}
	}

	fmt.Printf("Migrated %d tenants\n", len(tenants))
	return nil
}

// migrateApps 迁移应用数据
func (m *Migrator) migrateApps() error {
	fmt.Println("Migrating apps...")

	// 从 D1 获取应用数据
	apps, err := m.fetchD1Apps()
	if err != nil {
		return err
	}

	for _, a := range apps {
		app := models.App{
			ID:          a.ID,
			TenantID:    a.TenantID,
			AppName:     a.AppName,
			Description: a.Description,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}

		if err := m.DB.FirstOrCreate(&app, models.App{ID: a.ID}).Error; err != nil {
			return fmt.Errorf("failed to create app %d: %w", a.ID, err)
		}
	}

	fmt.Printf("Migrated %d apps\n", len(apps))
	return nil
}

// migrateMappings 迁移映射数据
func (m *Migrator) migrateMappings() error {
	fmt.Println("Migrating mappings...")

	// 从 D1 获取映射数据
	mappings, err := m.fetchD1Mappings()
	if err != nil {
		return err
	}

	for _, mData := range mappings {
		mapping := models.Mapping{
			ID:         mData.ID,
			TenantID:   mData.TenantID,
			AppID:      mData.AppID,
			PublicIP:   mData.PublicIP,
			PublicPort: mData.PublicPort,
			LocalIP:    mData.LocalIP,
			LocalPort:  mData.LocalPort,
			Protocol:   mData.Protocol,
			CreatedAt:  mData.CreatedAt,
			UpdatedAt:  mData.UpdatedAt,
		}

		if err := m.DB.FirstOrCreate(&mapping, models.Mapping{ID: mData.ID}).Error; err != nil {
			return fmt.Errorf("failed to create mapping %d: %w", mData.ID, err)
		}
	}

	fmt.Printf("Migrated %d mappings\n", len(mappings))
	return nil
}

// fetchD1Tenants 从 D1 获取租户数据
func (m *Migrator) fetchD1Tenants() ([]D1Record, error) {
	// 这里模拟从 D1 获取数据
	// 实际使用时，可以通过 Cloudflare API 或导出 SQL 文件
	// 简化版本：返回空列表，让用户手动导入
	fmt.Println("Note: Please export data from D1 manually and use ImportFromJSON method")
	return []D1Record{}, nil
}

// fetchD1Apps 从 D1 获取应用数据
func (m *Migrator) fetchD1Apps() ([]D1Record, error) {
	return []D1Record{}, nil
}

// fetchD1Mappings 从 D1 获取映射数据
func (m *Migrator) fetchD1Mappings() ([]D1Record, error) {
	return []D1Record{}, nil
}

// ImportFromJSON 从 JSON 文件导入数据
func (m *Migrator) ImportFromJSON(tenantsJSON, appsJSON, mappingsJSON string) error {
	fmt.Println("Importing data from JSON...")

	// 导入租户
	if tenantsJSON != "" {
		var tenants []models.Tenant
		if err := json.Unmarshal([]byte(tenantsJSON), &tenants); err != nil {
			return fmt.Errorf("failed to parse tenants JSON: %w", err)
		}

		for _, t := range tenants {
			if err := m.DB.FirstOrCreate(&t, models.Tenant{ID: t.ID}).Error; err != nil {
				return fmt.Errorf("failed to create tenant: %w", err)
			}
		}
		fmt.Printf("Imported %d tenants\n", len(tenants))
	}

	// 导入应用
	if appsJSON != "" {
		var apps []models.App
		if err := json.Unmarshal([]byte(appsJSON), &apps); err != nil {
			return fmt.Errorf("failed to parse apps JSON: %w", err)
		}

		for _, a := range apps {
			if err := m.DB.FirstOrCreate(&a, models.App{ID: a.ID}).Error; err != nil {
				return fmt.Errorf("failed to create app: %w", err)
			}
		}
		fmt.Printf("Imported %d apps\n", len(apps))
	}

	// 导入映射
	if mappingsJSON != "" {
		var mappings []models.Mapping
		if err := json.Unmarshal([]byte(mappingsJSON), &mappings); err != nil {
			return fmt.Errorf("failed to parse mappings JSON: %w", err)
		}

		for _, m := range mappings {
			if err := m.DB.FirstOrCreate(&m, models.Mapping{ID: m.ID}).Error; err != nil {
				return fmt.Errorf("failed to create mapping: %w", err)
			}
		}
		fmt.Printf("Imported %d mappings\n", len(mappings))
	}

	fmt.Println("Import completed!")
	return nil
}

// ExportToJSON 导出数据到 JSON
func (m *Migrator) ExportToJSON() (map[string]string, error) {
	result := make(map[string]string)

	// 导出租户
	var tenants []models.Tenant
	if err := m.DB.Find(&tenants).Error; err != nil {
		return nil, fmt.Errorf("failed to export tenants: %w", err)
	}
	tenantsJSON, _ := json.MarshalIndent(tenants, "", "  ")
	result["tenants"] = string(tenantsJSON)

	// 导出应用
	var apps []models.App
	if err := m.DB.Find(&apps).Error; err != nil {
		return nil, fmt.Errorf("failed to export apps: %w", err)
	}
	appsJSON, _ := json.MarshalIndent(apps, "", "  ")
	result["apps"] = string(appsJSON)

	// 导出映射
	var mappings []models.Mapping
	if err := m.DB.Find(&mappings).Error; err != nil {
		return nil, fmt.Errorf("failed to export mappings: %w", err)
	}
	mappingsJSON, _ := json.MarshalIndent(mappings, "", "  ")
	result["mappings"] = string(mappingsJSON)

	return result, nil
}
