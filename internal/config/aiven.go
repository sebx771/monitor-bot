package config

import (
	"os"
	"fmt"
	"errors"
)

type Credential struct {
	Token   string
	Project string
}

type AivenConfig struct {
	aivenCredentials []Credential
}

func NewAivenConfig() (*AivenConfig, error) {
	config := &AivenConfig{}

	if err := config.LoadAivenCredentials(); err != nil {
		return nil, err
	}

	return config, nil
}

func (a *AivenConfig) LoadAivenCredentials() error {
	credentials := []Credential{}

	for i := 1; ; i++ {
		token := os.Getenv(fmt.Sprintf("AIVEN_TOKEN_%d", i))
		project := os.Getenv(fmt.Sprintf("AIVEN_PROJECT_%d", i))

		if token == "" && project == "" {
			break
		}
	
		if token == "" || project == "" {
			return fmt.Errorf(
				"las variables AIVEN_TOKEN_%d y AIVEN_PROJECT_%d deben estar ambas definidas",
				i,
				i,
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
			credentials = append(credentials, Credential{
				Token:   token,
				Project: project,
			})
		}
	}

	if len(credentials) == 0 {
		return errors.New(
			"se requiere al menos una credencial AIVEN",
		)
	}

	a.aivenCredentials = credentials

	return nil
}