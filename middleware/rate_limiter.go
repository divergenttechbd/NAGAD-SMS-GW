package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	// "myproject/utils"
	"net/http"
	"time"
)

// RateLimiter middleware
func RateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// Increment the request count for this IP
		count, err := redisClient.Incr(c, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Set expiration for the key if this is the first request
		if count == 1 {
			redisClient.Expire(c, key, time.Minute)
		}

		// Allow a maximum of 100 requests per minute
		if count > 100 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}