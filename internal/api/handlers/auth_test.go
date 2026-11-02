package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginSetsCrossSiteSecureSessionCookie(t *testing.T) {
	previousMode := gin.Mode()
	gin.SetMode(gin.ReleaseMode)
	t.Cleanup(func() {
		gin.SetMode(previousMode)
	})

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	authService := auth.NewSessionService("admin", string(passwordHash))
	handler := NewAuthHandler(authService)
	router := gin.New()
	router.POST("/auth/login", handler.Login)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(`{"username":"admin","password":"password"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one session cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != sessionCookieName {
		t.Errorf("expected cookie %q, got %q", sessionCookieName, cookie.Name)
	}
	if cookie.SameSite != http.SameSiteNoneMode {
		t.Errorf("expected SameSite=None, got %v", cookie.SameSite)
	}
	if !cookie.Secure {
		t.Error("expected Secure cookie in release mode")
	}
	if !cookie.HttpOnly {
		t.Error("expected HttpOnly cookie")
	}
}
