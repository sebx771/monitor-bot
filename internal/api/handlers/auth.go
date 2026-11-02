package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/auth"
)

const (
	sessionCookieName = "monitor_session"
	sessionMaxAge     = int(8 * time.Hour / time.Second)
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthHandler manages admin authentication, session lifecycle, and status checks.
type AuthHandler struct {
	authSvc *auth.SessionService
}

func NewAuthHandler(authSvc *auth.SessionService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

func setSessionCookie(c *gin.Context, value string, maxAge int) {
	secure := gin.Mode() == gin.ReleaseMode
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(sessionCookieName, value, maxAge, "/", "", secure, true)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	sessionID, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	setSessionCookie(c, sessionID, sessionMaxAge)

	c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie(sessionCookieName)
	if err == nil && sessionID != "" {
		h.authSvc.Logout(sessionID)
	}

	setSessionCookie(c, "", -1)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	sessionID, err := c.Cookie(sessionCookieName)
	if err != nil || sessionID == "" || !h.authSvc.ValidateSession(sessionID) {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"username":      h.authSvc.GetAdminUsername(),
	})
}
