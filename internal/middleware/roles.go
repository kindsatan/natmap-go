package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
    allowed := map[string]struct{}{}
    for _, r := range roles { allowed[r] = struct{}{} }
    return func(c *gin.Context) {
        v, ok := c.Get("role")
        if !ok { c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return }
        role, _ := v.(string)
        if _, ok := allowed[role]; !ok { c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return }
        c.Next()
    }
}
