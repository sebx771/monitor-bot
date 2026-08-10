package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Credential struct {
	Token   string
	Project string
}

type Config struct {
	
	hostMC string
	portMC uint16

	serverId string
	storagePath string

	aivenCredentials []Credential
}

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := loadEnv(".env"); err != nil {
		return nil, err
	}

	if err := config.GetValues(); err != nil {
		return nil, err
	}

	if err := ensureStorageFile(config.storagePath); err != nil {
		return nil, err
	}

	return config, nil
}
func loadEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(value))
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

func (con *Config) GetValues() error{
	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		return err
	}

	host := os.Getenv("HOST")
	if host == ""{
	  return errors.New("La Variable Host es obligatoria")
	}

	server := os.Getenv("SERVER_ID")
	if server == ""{
		return errors.New("La Variable SERVER_ID es obligatoria")
	}

	path := os.Getenv("STORAGE_PATH")
	if path == ""{
	    return errors.New("La Variable STORAGE_PATH es obligatoria")
	}

	credentials, err := loadAivenCredentials()
	if err != nil {
		return err
	}

	con.hostMC = host
	con.portMC = uint16(port)
	con.serverId = server
	con.storagePath = path
	con.aivenCredentials = credentials

	return nil
}

// loadAivenCredentials lee las credenciales indexadas AIVEN_TOKEN_i /
// AIVEN_PROJECT_i hasta que no exista la variable con el índice siguiente.
// Si no hay credenciales indexadas, usa las legacy AIVEN_TOKEN/AIVEN_PROJECT.
func loadAivenCredentials() ([]Credential, error) {
	credentials := []Credential{}

	for i := 1; ; i++ {
		token := os.Getenv(fmt.Sprintf("AIVEN_TOKEN_%d", i))
		project := os.Getenv(fmt.Sprintf("AIVEN_PROJECT_%d", i))

		if token == "" && project == "" {
			break
		}

		if token == "" || project == "" {
			return nil, fmt.Errorf(
				"Las Variables AIVEN_TOKEN_%d y AIVEN_PROJECT_%d deben estar ambas definidas",
				i, i,
			)
		}

		credentials = append(credentials, Credential{
			Token:   token,
			Project: project,
		})
	}

	if len(credentials) == 0 {
		token := os.Getenv("AIVEN_TOKEN")
		project := os.Getenv("AIVEN_PROJECT")

		if token != "" && project != "" {
			return []Credential{{Token: token, Project: project}}, nil
		}

		return nil, errors.New(
			"Se requiere al menos una credencial AIVEN (AIVEN_TOKEN_1/AIVEN_PROJECT_1 " +
				"o las legacy AIVEN_TOKEN/AIVEN_PROJECT)",
		)
	}

	return credentials, nil
}

//getters
func (con *Config) GetHostMC() string {
	return con.hostMC
}

func (con *Config) GetPortMC() uint16 {
	return con.portMC
}

func (con *Config) GetServerID() string {
	return con.serverId
}

func (con *Config) GetStoragePath() string {
	return con.storagePath
}

func (con *Config) GetAivenCredentials() []Credential{
    return con.aivenCredentials
}


