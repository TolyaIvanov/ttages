package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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
	}
}
