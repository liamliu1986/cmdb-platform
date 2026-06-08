package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"cmdb-api/pkg/response"
)

type clientLimiter struct {
	count    int
	window   time.Time
}

var (
	limiters = make(map[string]*clientLimiter)
	limiterMu sync.RWMutex
)

// RateLimit limits requests per IP to maxRequests per windowDuration
func RateLimit(maxRequests int, windowDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiterMu.Lock()
		now := time.Now()
		cl, exists := limiters[ip]
		if !exists || now.Sub(cl.window) > windowDuration {
			limiters[ip] = &clientLimiter{count: 1, window: now}
			cl = limiters[ip]
		} else {
			cl.count++
		}
		count := cl.count
		limiterMu.Unlock()

		if count > maxRequests {
			response.ErrorWithStatus(c, http.StatusTooManyRequests, 10020, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
