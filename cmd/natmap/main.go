package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"natmap/internal/cache"
	"natmap/internal/config"
	"natmap/internal/handlers"
	"natmap/internal/middleware"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	migrate := flag.Bool("migrate", false, "Run database migration")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger := initLogger(cfg.Log)
	defer logger.Sync()

	logger.Info("Starting NATMap server",
		zap.String("version", "1.0.0"),
		zap.Int("port", cfg.Server.Port),
	)

	// 连接数据库
	db, err := handlers.NewDatabase(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// 自动迁移
	if *migrate {
		logger.Info("Running database migration...")
		if err := db.AutoMigrate(); err != nil {
			logger.Fatal("Failed to migrate database", zap.Error(err))
		}
		logger.Info("Database migration completed")
		return
	}

	// 初始化缓存
	var appCache cache.Cache
	if cfg.Cache.Enabled {
		appCache = cache.NewMemoryCache(
			time.Duration(cfg.Cache.TTL)*time.Second,
			time.Duration(cfg.Cache.CleanupInterval)*time.Second,
		)
		logger.Info("Cache enabled", zap.Int("ttl_seconds", cfg.Cache.TTL))
	} else {
		appCache = cache.NewNoOpCache()
		logger.Info("Cache disabled")
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	// 创建处理器
	mappingHandler := handlers.NewMappingHandler(db, appCache, logger)
	adminHandler := handlers.NewAdminHandler(db, appCache, logger)

	// 公开 API（无需认证）
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "NATMap API is running",
			})
		})

		// 获取映射（公开接口）
		api.GET("/get", mappingHandler.GetMapping)

		// 更新映射（需要 token 认证）
		// 简化版：暂时不使用中间件，由处理器内部处理
		api.POST("/update", mappingHandler.UpdateMapping)
	}

	// 管理后台 API（需要 Basic Auth）
	// 注意: 生产环境请启用认证
	admin := api.Group("/admin")
	// admin.Use(middleware.BasicAuth(db.Conn))
	{
		// 租户管理
		admin.GET("", func(c *gin.Context) {
			typeParam := c.Query("type")
			switch typeParam {
			case "tenant":
				adminHandler.ListTenants(c)
			case "app":
				adminHandler.ListApps(c)
			case "mapping":
				adminHandler.ListMappings(c)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
			}
		})

		// 创建资源
		admin.POST("", func(c *gin.Context) {
			typeParam := c.Query("type")
			switch typeParam {
			case "tenant":
				adminHandler.CreateTenant(c)
			case "app":
				adminHandler.CreateApp(c)
			case "mapping":
				adminHandler.CreateMapping(c)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
			}
		})

		// 更新资源
		admin.PUT("", func(c *gin.Context) {
			typeParam := c.Query("type")
			switch typeParam {
			case "tenant":
				adminHandler.UpdateTenant(c)
			case "app":
				adminHandler.UpdateApp(c)
			case "mapping":
				adminHandler.UpdateMapping(c)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
			}
		})

		// 删除资源
		admin.DELETE("", func(c *gin.Context) {
			typeParam := c.Query("type")
			switch typeParam {
			case "tenant":
				adminHandler.DeleteTenant(c)
			case "app":
				adminHandler.DeleteApp(c)
			case "mapping":
				adminHandler.DeleteMapping(c)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
			}
		})
	}

	// 静态文件服务（React 管理后台）
	r.Static("/assets", "./web/assets")
	r.StaticFile("/favicon.svg", "./web/favicon.svg")
	r.StaticFile("/icons.svg", "./web/icons.svg")

	// 所有其他路由都返回 index.html（支持 React Router）
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.Int("port", cfg.Server.Port))

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	if err := srv.Close(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// initLogger 初始化日志
func initLogger(cfg config.LogConfig) *zap.Logger {
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core, zap.AddCaller())
}
