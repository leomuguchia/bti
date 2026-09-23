package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const bearerPrefix = "Bearer "

func APIKeyAuth(expectedKey string) gin.HandlerFunc {
	// Fail closed. An unset key must reject everything, not
	// accidentally accept requests that also have no header.
	if expectedKey == "" {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server misconfigured"})
		}
	}

	expected := []byte(expectedKey)
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, bearerPrefix) {
			c.Header("WWW-Authenticate", `Bearer realm="banking"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		got := []byte(strings.TrimPrefix(header, bearerPrefix))
		if subtle.ConstantTimeCompare(got, expected) != 1 {
			c.Header("WWW-Authenticate", `Bearer realm="banking"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
