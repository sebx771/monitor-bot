package auth

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// newTestService is a helper that creates a SessionService with known credentials.
func newTestService(t *testing.T, duration ...time.Duration) (*SessionService, string) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to generate bcrypt hash: %v", err)
	}
	svc := &SessionService{
		sessions:          make(map[string]Session),
		adminUsername:     "admin",
		adminPasswordHash: string(hash),
		sessionDuration:   8 * time.Hour,
	}
	if len(duration) > 0 {
		svc.sessionDuration = duration[0]
	}
	return svc, "testpass"
}

func TestLogin_Success(t *testing.T) {
	svc, pass := newTestService(t)

	sessionID, err := svc.Login("admin", pass)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if sessionID == "" {
		t.Fatal("expected non-empty sessionID")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.Login("admin", "wrongpassword")

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_WrongUsername(t *testing.T) {
	svc, pass := newTestService(t)

	_, err := svc.Login("notadmin", pass)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_SessionIDNeverEmpty(t *testing.T) {
	svc, pass := newTestService(t)

	for i := range 5 {
		id, err := svc.Login("admin", pass)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		if id == "" {
			t.Fatalf("iteration %d: sessionID must not be empty", i)
		}
	}
}

func TestValidateSession_Valid(t *testing.T) {
	svc, pass := newTestService(t)

	sessionID, err := svc.Login("admin", pass)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if !svc.ValidateSession(sessionID) {
		t.Fatal("expected session to be valid right after creation")
	}
}

func TestValidateSession_NotFound(t *testing.T) {
	svc, _ := newTestService(t)

	if svc.ValidateSession("nonexistent-session-id-xyz") {
		t.Fatal("expected false for unknown session ID")
	}
}

func TestValidateSession_Expired(t *testing.T) {
	// Use a negative duration so sessions are immediately expired
	svc, pass := newTestService(t, -1*time.Second)

	sessionID, err := svc.Login("admin", pass)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Session was created with ExpiresAt in the past
	if svc.ValidateSession(sessionID) {
		t.Fatal("expected expired session to be invalid")
	}
}

func TestLogout_InvalidatesSession(t *testing.T) {
	svc, pass := newTestService(t)

	sessionID, err := svc.Login("admin", pass)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	svc.Logout(sessionID)

	if svc.ValidateSession(sessionID) {
		t.Fatal("expected session to be invalid after logout")
	}
}

func TestNewSessionService_PanicOnEmptyUsername(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty username, got none")
		}
	}()
	NewSessionService("", "$2a$04$validhash")
}

func TestNewSessionService_PanicOnEmptyPasswordHash(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty passwordHash, got none")
		}
	}()
	NewSessionService("admin", "")
}
