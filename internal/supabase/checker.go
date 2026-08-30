package supabase

import (
	"fmt"

	"github.com/sebx771/monitor-bot/internal/logger"
)

type Checker struct {
	client *Client
	logger *logger.Logger
}

func NewChecker(client *Client, logger *logger.Logger) *Checker {
	return &Checker{
		client: client,
		logger: logger,
	}
}

func (c *Checker) Check() error {
	projects, err := c.client.ListProjects()
	if err != nil {
		return fmt.Errorf("error obteniendo proyectos de Supabase: %w", err)
	}

	for _, project := range projects {
		if project.Status != "INACTIVE" {
			continue
		}

		c.logger.Info(
			"proyecto de Supabase inactivo, intentando restaurar",
			"project", project.Name,
			"ref", project.Ref,
		)

		if err := c.client.RestoreProject(project.Ref); err != nil {
			c.logger.Error(
				"error restaurando proyecto de Supabase",
				"project", project.Name,
				"ref", project.Ref,
				"error", err,
			)

			continue
		}

		c.logger.Info(
			"proyecto de Supabase restaurado",
			"project", project.Name,
			"ref", project.Ref,
		)
	}

	return nil
}