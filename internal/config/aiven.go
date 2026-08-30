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
	enabled bool
	aivenCredentials []Credential
}

func NewAivenConfig() (*AivenConfig, error) {
	config := &AivenConfig{}

	isEnabled, err := GetEnabledVar("AIVEN_ENABLED")
	if err != nil {
		return nil, err
	}

	config.enabled = isEnabled

	if !isEnabled {
		return config, nil
	}

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
				"the AIVEN_TOKEN_%d and AIVEN_PROJECT_%d variables must both be defined",
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
			"at least one AIVEN credential is required",
		)
	}

	a.aivenCredentials = credentials

	return nil
}

func (a *AivenConfig) IsEnabled() bool {
	return a.enabled
}

func (a *AivenConfig) GetCredentials() []Credential {
	return a.aivenCredentials
}