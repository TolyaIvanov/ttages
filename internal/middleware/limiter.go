package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mutex    sync.Mutex
}

func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (rl *RateLimiter) LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.mutex.Lock()
		if _, exists := rl.limiters[ip]; !exists {
			rl.limiters[ip] = rate.NewLimiter(rate.Limit(10), 10)
		}
		limiter := rl.limiters[ip]
		rl.mutex.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
