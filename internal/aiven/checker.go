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

// Check revisa los servicios del proyecto y enciende los apagados.
// Un error en un servicio no detiene el chequeo del resto.
func (c *AivenChecker) Check() error {
	services, err := c.client.GetServices(c.project)
	if err != nil {
		return fmt.Errorf("error obteniendo Aiven services (proyecto %s): %w", c.project, err)
	}

	for _, service := range services {
		c.log.Info("servicio revisado", "proyecto", c.project, "servicio", service.Name, "estado", service.State)

		if service.State != "POWEROFF" {
			continue
		}

		c.log.Info("iniciando servicio", "proyecto", c.project, "servicio", service.Name)

		if err := c.client.StartService(c.project, service); err != nil {
			c.log.Error(
				"error iniciando servicio",
				"proyecto", c.project,
				"servicio", service.Name,
				"error", err,
			)
			continue
		}
	}

	return nil
}
