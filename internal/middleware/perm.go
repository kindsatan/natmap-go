package middleware

import (
    "net/http"

    "natmap-go/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RequirePermission(db *gorm.DB, resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        v, ok := c.Get("role")
        if !ok { c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return }
        role, _ := v.(string)
        var p models.Permission
        if err := db.Where("role = ? AND resource = ? AND action = ? AND allowed = ?", role, resource, action, true).First(&p).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
            return
        }
        c.Next()
    }
}
