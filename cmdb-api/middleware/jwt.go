package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cmdb-api/config"
	"cmdb-api/pkg/jwtutil"
	"cmdb-api/pkg/response"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 10010, "missing authorization header")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 10011, "invalid authorization header format")
			c.Abort()
			return
		}
		claims, err := jwtutil.ParseToken(parts[1], cfg.JWTSecret)
		if err != nil {
			response.ErrorWithStatus(c, http.StatusUnauthorized, 10012, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
