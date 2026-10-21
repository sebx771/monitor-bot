package middlewares

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuth valida que la cabecera X-API-Key coincida con expectedKey
// utilizando comparación en tiempo constante (constant-time).
func ApiKeyAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientKey := c.GetHeader("X-API-Key")

		if expectedKey == "" || clientKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}

		if subtle.ConstantTimeCompare([]byte(clientKey), []byte(expectedKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			return
		}

		c.Next()
	}
}
