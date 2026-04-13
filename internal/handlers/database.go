package handlers

import (
	"time"

	"natmap/internal/config"
	"natmap/internal/models"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 数据库连接
type DB struct {
	Conn   *gorm.DB
	Logger *zap.Logger
}

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.DatabaseConfig, log *zap.Logger) (*DB, error) {
	// 配置 GORM 日志
	gormLogger := logger.Default
	if log != nil {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return &DB{
		Conn:   db,
		Logger: log,
	}, nil
}

// AutoMigrate 自动迁移数据库表
func (db *DB) AutoMigrate() error {
	return db.Conn.AutoMigrate(
		&models.Tenant{},
		&models.App{},
		&models.Mapping{},
		&models.User{},
	)
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	sqlDB, err := db.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
