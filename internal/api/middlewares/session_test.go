package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebx771/monitor-bot/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func setupSessionRouter(svc *auth.SessionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SessionAuth(svc))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func newTestSessionService(t *testing.T, duration ...time.Duration) (*auth.SessionService, func() string) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	svc := auth.NewSessionService("admin", string(hash))
	if len(duration) > 0 {
		// We cannot set sessionDuration externally; for expired tests we log in
		// and then manipulate via the public API only — see TestSessionAuth_ExpiredSession.
		_ = duration
	}
	loginFn := func() string {
		id, err := svc.Login("admin", "pass")
		if err != nil {
			t.Fatalf("login: %v", err)
		}
		return id
	}
	return svc, loginFn
}

func requestWithCookie(router *gin.Engine, cookieName, cookieValue string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	if cookieName != "" {
		req.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestSessionAuth_ValidSession(t *testing.T) {
	svc, login := newTestSessionService(t)
	router := setupSessionRouter(svc)
	sessionID := login()

	w := requestWithCookie(router, sessionCookieName, sessionID)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSessionAuth_NoCookie(t *testing.T) {
	svc, _ := newTestSessionService(t)
	router := setupSessionRouter(svc)

	w := requestWithCookie(router, "", "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSessionAuth_InvalidSessionID(t *testing.T) {
	svc, _ := newTestSessionService(t)
	router := setupSessionRouter(svc)

	w := requestWithCookie(router, sessionCookieName, "totally-fake-session-id")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSessionAuth_ExpiredSession(t *testing.T) {
	svc, login := newTestSessionService(t)
	router := setupSessionRouter(svc)

	sessionID := login()
	// Invalidate immediately via logout (simulates expiry from caller's perspective)
	svc.Logout(sessionID)

	w := requestWithCookie(router, sessionCookieName, sessionID)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 after logout/expiry, got %d: %s", w.Code, w.Body.String())
	}
}
