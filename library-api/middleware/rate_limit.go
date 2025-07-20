package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"github.com/example/library-api/config"
)

// IPRateLimiter 基于IP的限流结构
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	rps int
	burst int
}

var limiter *IPRateLimiter

// InitRateLimiter 初始化限流中间件
func InitRateLimiter(cfg *config.Config) {
	limiter = &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		rps: cfg.RateLimitRPS,
		burst: cfg.RateLimitBurst,
	}
}

// addIP 如果IP不存在则添加到限流映射
func (i *IPRateLimiter) addIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	lim, exists := i.ips[ip]
	if !exists {
		lim = rate.NewLimiter(rate.Limit(i.rps), i.burst)
		i.ips[ip] = lim
	}

	return lim
}

// getLimiter 获取IP对应的限流器
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	lim, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.addIP(ip)
	}
	return lim
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		lim := limiter.getLimiter(ip)

		if !lim.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}