package config

import (
	"errors"
	"os"
	"strings"
)

type GlobalConfig struct {
	secretKey         string
	adminUserName     string
	adminPasswordHash string
}

func NewGlobalConfig() (*GlobalConfig, error) {
	adminUserName := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	if adminUserName == "" {
		return nil, errors.New("the ADMIN_USERNAME variable is required")
	}

	adminPasswordHash := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD_HASH"))
	if adminPasswordHash == "" {
		return nil, errors.New("the ADMIN_PASSWORD_HASH variable is required")
	}

	return &GlobalConfig{
		secretKey:         strings.TrimSpace(os.Getenv("SECRET_KEY")),
		adminUserName:     adminUserName,
		adminPasswordHash: adminPasswordHash,
	}, nil
}

func (gc *GlobalConfig) GetSecretKey() string {
	return gc.secretKey
}

func (gc *GlobalConfig) GetAdminUserName() string {
	return gc.adminUserName
}

func (gc *GlobalConfig) GetAdminPasswordHash() string {
	return gc.adminPasswordHash
}
