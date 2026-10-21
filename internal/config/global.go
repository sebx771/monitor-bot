package config


import(
	"os"
)


type GlobalConfig struct {
	secretKey string
}

func NewGlobalConfig() (*GlobalConfig, error) {
	return &GlobalConfig{
		secretKey: os.Getenv("SECRET_KEY"),
	}, nil
}

func (gc *GlobalConfig) GetSecretKey() string {
	return gc.secretKey
}
