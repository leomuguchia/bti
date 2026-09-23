package middleware

import "github.com/gin-gonic/gin"

// No-op placeholder for the demo. Production would use a token bucket keyed
// on API key / IP (golang.org/x/time/rate), backed by Redis if multi-instance.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
