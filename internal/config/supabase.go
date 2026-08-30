package config

import (
	"errors"
	"os"
)

type SupabaseConfig struct {
	Enabled bool
	Token   string
}

func NewSupabaseConfig() (*SupabaseConfig, error) {
	config := &SupabaseConfig{}

	isEnabled, err := GetEnabledVar("SUPABASE_ENABLED")
	if err != nil {
		return nil, err
	}

	config.Enabled = isEnabled

	if !config.Enabled {
		return config, nil
	}

	if err := config.LoadSupabaseCredentials(); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *SupabaseConfig) LoadSupabaseCredentials() error {
	token := os.Getenv("SUPABASE_TOKEN")

	if token == "" {
		return errors.New("se requiere SUPABASE_TOKEN")
	}

	s.Token = token

	return nil
}

func (s *SupabaseConfig) GetEnabled() bool {
	return s.Enabled
}

func (s *SupabaseConfig) GetToken() string {
	return s.Token
}