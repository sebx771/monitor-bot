package config

import (
	"os"
	"strings"

	"github.com/sebx771/monitor-bot/internal/logger"
)


type Config struct {
	globalConfig  struct{} // de momento lo dejo asi , hasta que haga los struct reales de las config
	aivenConfig   AivenConfig
	aternosConfig  AternosConfig
}

func NewConfig() (*Config, error){
	logs:= logger.NewLogger("CONFIG")
	config := &Config{}

	logs.Debug("Iniciando configuracion del sistema")

	if err := loadEnv(".env"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// en este bloque ira la declaracion de las demas config para posteriormente instancear
	//  la config raiz 
	aiven,err:= NewAivenConfig()
	if err!= nil {
		return nil , err
	}

	aternos , err := NewAternosConfig()
	if err != nil {
		return nil , err
	}

	
    config.aivenConfig= *aiven
	config.aternosConfig= *aternos
	

    return  config, nil

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
