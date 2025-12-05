package main

import (
    "log"
    "net/http"
    "time"

    "natmap-go/internal/auth"
    "natmap-go/internal/config"
    "natmap-go/internal/db"
    "natmap-go/internal/handlers"
    "natmap-go/internal/middleware"
    "natmap-go/internal/models"
    _ "natmap-go/internal/docs"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
    cfg := config.Load()
    d, err := db.SetupDB(cfg.SQLitePath)
    if err != nil {
        log.Fatalf("db error: %v", err)
    }
    if err := d.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Permission{}); err != nil {
        log.Fatalf("migrate error: %v", err)
    }
    if cfg.SeedUser {
        var count int64
        d.Model(&models.User{}).Count(&count)
        if count == 0 && cfg.SeedUsername != "" && cfg.SeedPassword != "" {
            hash, _ := auth.HashPassword(cfg.SeedPassword, cfg.BcryptCost)
            u := models.User{Username: cfg.SeedUsername, PasswordHash: hash, Role: "admin", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
            d.Create(&u)
        }
        if cfg.SeedUsername != "" {
            var su models.User
            if err := d.Where("username = ?", cfg.SeedUsername).First(&su).Error; err == nil {
                d.Model(&su).Update("role", "admin")
            }
        }
    }

    r := gin.New()
    r.Use(gin.Recovery())
    r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    ah := handlers.NewAuthHandler(d, cfg.JWTSecret, cfg.TokenTTL, cfg.RefreshTTL, cfg.BcryptCost)
    uh := handlers.NewUserHandler(d)
    admin := handlers.NewAdminHandler(d)

    v1 := r.Group("/api/v1")
    v1.POST("/auth/login", ah.Login)
    v1.POST("/auth/register", ah.Register)
    v1.POST("/auth/refresh", ah.Refresh)
    v1.POST("/auth/logout", middleware.AuthMiddleware(cfg.JWTSecret), ah.Logout)
    v1.POST("/auth/logout_all", middleware.AuthMiddleware(cfg.JWTSecret), ah.LogoutAll)

    protected := v1.Group("")
    protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
    protected.GET("/me", uh.Me)
    protectedAdmin := protected.Group("/admin")
    protectedAdmin.Use(middleware.RequireRole("admin"))
    protectedAdmin.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
    protectedAdmin.GET("/users", admin.ListUsers)
    protectedAdmin.PUT("/users/:id/role", admin.SetUserRole)
    protectedAdmin.GET("/permissions", admin.ListPermissions)
    protectedAdmin.POST("/permissions", admin.CreatePermission)
    protectedAdmin.DELETE("/permissions/:id", admin.DeletePermission)

    protected.GET("/reports", middleware.RequirePermission(d, "reports", "read"), func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"items": []string{"r1", "r2"}})
    })

    if err := r.Run(cfg.HTTPAddr); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
