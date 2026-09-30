package supabase

import (
	"fmt"
	"strings"

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

// isInactive checks if the project status indicates it is suspended (INACTIVE or PAUSED).
func isInactive(status string) bool {
	s := strings.ToUpper(strings.TrimSpace(status))
	return s == "INACTIVE" || s == "PAUSED"
}

func (c *Checker) Check() error {
	projects, err := c.client.ListProjects()
	if err != nil {
		return fmt.Errorf("error getting Supabase projects: %w", err)
	}

	for _, project := range projects {
		if !isInactive(project.Status) {
			c.logger.Info(
				"Supabase project active",
				"project", project.Name,
				"ref", project.Ref,
				"status", project.Status,
			)
			continue
		}

		c.logger.Info(
			"Supabase project inactive or paused, attempting to restore",
			"project", project.Name,
			"ref", project.Ref,
			"status", project.Status,
		) 

		if err := c.client.RestoreProject(project.Ref); err != nil {
			c.logger.Error(
				"error restoring Supabase project",
				"project", project.Name,
				"ref", project.Ref,
				"error", err,
			)

			continue
		}

		c.logger.Info(
			"Supabase project restored",
			"project", project.Name,
			"ref", project.Ref,
		)
	}

	return nil
}