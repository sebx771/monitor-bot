package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/auth"
)

const sessionCookieName = "monitor_session"

// SessionAuth validates the "monitor_session" HttpOnly cookie against the
// SessionService. Responds 401 and aborts if the session is missing or invalid.
func SessionAuth(svc *auth.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(sessionCookieName)
		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		if !svc.ValidateSession(sessionID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		c.Next()
	}
}
