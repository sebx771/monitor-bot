package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

type AternosConfig struct {
	enabled     bool
	host        string
	port        uint16
	serverId    string
	storagePath string

	headless bool

	githubToken string
	gistId      string
}

func NewAternosConfig() (*AternosConfig, error) {
	config := &AternosConfig{}

	isEnabled,err:= GetEnabledVar("ATERNOS_ENABLED")
	if err != nil {
		return nil, err
	}
	config.enabled= isEnabled
	if !config.enabled {
      return config , nil
	}
	

	if err := config.GetValues(); err != nil {
		return nil, err
	}

	return config, nil
}

func (At *AternosConfig) GetValues() error {

	
	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		return err
	}

	host := os.Getenv("HOST")
	if host == "" {
		return errors.New("La Variable HOST es obligatoria")
	}

	server := os.Getenv("SERVER_ID")
	if server == "" {
		return errors.New("La Variable SERVER_ID es obligatoria")
	}

	path := os.Getenv("STORAGE_PATH")
	if path == "" {
		return errors.New("La Variable STORAGE_PATH es obligatoria")
	}

	headless := os.Getenv("HEADLESS")
	if headless == "" {
		return errors.New("La Variable HEADLESS es obligatoria")
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return errors.New("La Variable GITHUB_TOKEN es obligatoria")
	}

	gistId := os.Getenv("GIST_ID")
	if gistId == "" {
		return errors.New("La Variable GIST_ID es obligatoria")
	}

	At.host = host
	At.port = uint16(port)
	At.serverId = server
	At.storagePath = path
	At.headless = headless == "true"
	At.githubToken = githubToken
	At.gistId = gistId

	if err := ensureStorageFile(At.storagePath); err != nil {
		return err
	}

	return nil
}

func ensureStorageFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte("{}"), 0o644)
}

func (At *AternosConfig) IsEnabled() bool {
	return At.enabled
}

func (At *AternosConfig) GetHost() string {
	return At.host
}

func (At *AternosConfig) GetPort() uint16 {
	return At.port
}

func (At *AternosConfig) GetServerID() string {
	return At.serverId
}

func (At *AternosConfig) GetStoragePath() string {
	return At.storagePath
}

func (At *AternosConfig) IsHeadless() bool {
	return At.headless
}

func (At *AternosConfig) GetGithubToken() string {
	return At.githubToken
}

func (At *AternosConfig) GetGistID() string {
	return At.gistId
}