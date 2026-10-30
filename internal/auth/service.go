package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is returned when username or password do not match.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Session holds metadata for an active admin session.
type Session struct {
	ID        string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionService manages in-memory admin sessions.
// It is safe for concurrent use.
type SessionService struct {
	mu                sync.RWMutex
	sessions          map[string]Session
	adminUsername     string
	adminPasswordHash string
	sessionDuration   time.Duration
}

// NewSessionService creates a new SessionService.
// Panics if username or passwordHash are empty.
func NewSessionService(username, passwordHash string) *SessionService {
	if username == "" {
		panic("auth: NewSessionService requires a non-empty username")
	}
	if passwordHash == "" {
		panic("auth: NewSessionService requires a non-empty passwordHash")
	}
	return &SessionService{
		sessions:          make(map[string]Session),
		adminUsername:     username,
		adminPasswordHash: passwordHash,
		sessionDuration:   8 * time.Hour,
	}
}

// Login validates credentials and returns a new session ID on success.
// Returns ErrInvalidCredentials if username or password are wrong.
// The session ID is never logged.
func (s *SessionService) Login(username, password string) (string, error) {
	usernameMatch := subtle.ConstantTimeCompare(
		[]byte(username),
		[]byte(s.adminUsername),
	) == 1

	passwordErr := bcrypt.CompareHashAndPassword(
		[]byte(s.adminPasswordHash),
		[]byte(password),
	)

	// Evaluate both to avoid timing side-channels on username
	if !usernameMatch || passwordErr != nil {
		return "", ErrInvalidCredentials
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	now := time.Now()
	s.mu.Lock()
	s.sessions[sessionID] = Session{
		ID:        sessionID,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionDuration),
	}
	s.mu.Unlock()

	return sessionID, nil
}

// Logout removes the session from the store, invalidating it immediately.
func (s *SessionService) Logout(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

// ValidateSession returns true if the session exists and has not expired.
// Expired sessions are lazily removed on access.
func (s *SessionService) ValidateSession(sessionID string) bool {
	s.mu.RLock()
	session, ok := s.sessions[sessionID]
	s.mu.RUnlock()

	if !ok {
		return false
	}

	if time.Now().After(session.ExpiresAt) {
		s.mu.Lock()
		delete(s.sessions, sessionID)
		s.mu.Unlock()
		return false
	}

	return true
}

// GetAdminUsername returns the configured admin username.
// Safe to call from handlers; does not expose the password hash.
func (s *SessionService) GetAdminUsername() string {
	return s.adminUsername
}

// generateSessionID creates a cryptographically secure random 32-byte hex string.
func generateSessionID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
