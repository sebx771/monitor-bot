package config

import (
	"os"
	"strings"
	"testing"
)

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	prev, exists := os.LookupEnv(key)
	os.Unsetenv(key)
	t.Cleanup(func() {
		if exists {
			os.Setenv(key, prev)
		} else {
			os.Unsetenv(key)
		}
	})
}

func TestNewGlobalConfig_ValidCredentials(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "admin_user")
	t.Setenv("ADMIN_PASSWORD_HASH", "$2a$12$e8kZ1qUuQj5v7pQ...")
	t.Setenv("SECRET_KEY", "super-secret-key")

	cfg, err := NewGlobalConfig()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected non-nil GlobalConfig")
	}

	if cfg.GetAdminUserName() != "admin_user" {
		t.Errorf("expected adminUserName to be %q, got %q", "admin_user", cfg.GetAdminUserName())
	}

	if cfg.GetAdminPasswordHash() != "$2a$12$e8kZ1qUuQj5v7pQ..." {
		t.Errorf("expected adminPasswordHash to be %q, got %q", "$2a$12$e8kZ1qUuQj5v7pQ...", cfg.GetAdminPasswordHash())
	}

	if cfg.GetSecretKey() != "super-secret-key" {
		t.Errorf("expected secretKey to be %q, got %q", "super-secret-key", cfg.GetSecretKey())
	}
}

func TestNewGlobalConfig_MissingOrEmptyAdminUsername(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD_HASH", "valid_hash")

	tests := []struct {
		name     string
		setupEnv func(t *testing.T)
	}{
		{
			name: "missing/unset variable",
			setupEnv: func(t *testing.T) {
				unsetEnv(t, "ADMIN_USERNAME")
			},
		},
		{
			name: "empty string",
			setupEnv: func(t *testing.T) {
				t.Setenv("ADMIN_USERNAME", "")
			},
		},
		{
			name: "whitespace only",
			setupEnv: func(t *testing.T) {
				t.Setenv("ADMIN_USERNAME", "   ")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)

			cfg, err := NewGlobalConfig()
			if err == nil {
				t.Fatalf("expected error for %s, but got nil", tt.name)
			}

			if cfg != nil {
				t.Fatalf("expected nil config on error, got: %+v", cfg)
			}

			expectedErr := "the ADMIN_USERNAME variable is required"
			if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("expected error containing %q, got: %q", expectedErr, err.Error())
			}
		})
	}
}

func TestNewGlobalConfig_MissingOrEmptyAdminPasswordHash(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "valid_admin")

	tests := []struct {
		name     string
		setupEnv func(t *testing.T)
	}{
		{
			name: "missing/unset variable",
			setupEnv: func(t *testing.T) {
				unsetEnv(t, "ADMIN_PASSWORD_HASH")
			},
		},
		{
			name: "empty string",
			setupEnv: func(t *testing.T) {
				t.Setenv("ADMIN_PASSWORD_HASH", "")
			},
		},
		{
			name: "whitespace only",
			setupEnv: func(t *testing.T) {
				t.Setenv("ADMIN_PASSWORD_HASH", "   ")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)

			cfg, err := NewGlobalConfig()
			if err == nil {
				t.Fatalf("expected error for %s, but got nil", tt.name)
			}

			if cfg != nil {
				t.Fatalf("expected nil config on error, got: %+v", cfg)
			}

			expectedErr := "the ADMIN_PASSWORD_HASH variable is required"
			if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("expected error containing %q, got: %q", expectedErr, err.Error())
			}
		})
	}
}
