package config

import (
	"os"
	"strings"

	"github.com/sebx771/monitor-bot/internal/logger"
)


type Config struct {
	globalConfig  GlobalConfig 
	aivenConfig   AivenConfig
	aternosConfig  AternosConfig
	supaBaseConfig SupabaseConfig
}

func NewConfig() (*Config, error){
	logs:= logger.NewLogger("CONFIG")
	config := &Config{}

	logs.Debug("Starting system configuration")
	if err := loadEnv(".env"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	global , err:= NewGlobalConfig()
	if err!= nil {
		return nil , err
	}

	aiven,err:= NewAivenConfig()
	if err!= nil {
		return nil , err
	}

	aternos , err := NewAternosConfig()
	if err != nil {
		return nil , err
	}
	supabase , err := NewSupabaseConfig()
	if err != nil {
		return  nil , err 
	}

	config.globalConfig= *global
    config.aivenConfig= *aiven
	config.aternosConfig= *aternos
	config.supaBaseConfig= *supabase
	

    return  config, nil

}

func (c *Config) GetAivenConfig() *AivenConfig {
	return &c.aivenConfig
}

func (c *Config) GetAternosConfig() *AternosConfig {
	return &c.aternosConfig
}
func (c *Config) GetSupaBaseConfig() *SupabaseConfig{
	return &c.supaBaseConfig
}

func (c *Config) GetGlobalConfig() *GlobalConfig {
	return &c.globalConfig
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
