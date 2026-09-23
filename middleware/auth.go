package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(expectedKey string) gin.HandlerFunc {
	expected := []byte(expectedKey)
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}
		got := []byte(c.GetHeader("X-API-Key"))
		// Constant-time compare so response timing can't leak the key.
		if subtle.ConstantTimeCompare(got, expected) != 1 {
			c.Header("WWW-Authenticate", `ApiKey realm="banking"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
