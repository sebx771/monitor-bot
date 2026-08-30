package aiven

import (
	"fmt"

	"github.com/sebx771/monitor-bot/internal/logger"
)

type AivenChecker struct {
	client  *Client
	project string
	log     *logger.Logger
}

func NewChecker(client *Client, project string) *AivenChecker {
	return &AivenChecker{
		client:  client,
		project: project,
		log:     logger.NewLogger("AIVEN"),
	}
}

// Check reviews the project's services and starts the ones that are off.
// An error in one service does not stop checking the rest.
func (c *AivenChecker) Check() error {
	services, err := c.client.GetServices(c.project)
	if err != nil {
		return fmt.Errorf("error getting Aiven services (project %s): %w", c.project, err)
	}

	for _, service := range services {
		c.log.Info("service checked", "project", c.project, "service", service.Name, "state", service.State)

		if service.State != "POWEROFF" {
			continue
		}

		c.log.Info("starting service", "project", c.project, "service", service.Name)

		if err := c.client.StartService(c.project, service); err != nil {
			c.log.Error(
				"error starting service",
				"project", c.project,
				"service", service.Name,
				"error", err,
			)
			continue
		}
	}

	return nil
}
