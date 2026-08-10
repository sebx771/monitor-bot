package aiven

import (
	"fmt"
	"log"
)

type AivenChecker struct {
	client  *Client
	project string
}

func NewChecker(client *Client, project string) *AivenChecker {
	return &AivenChecker{
		client:  client,
		project: project,
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
		log.Printf("[Aiven/%s] Service: %s | State: %s", c.project, service.Name, service.State)

		if service.State != "POWEROFF" {
			continue
		}

		log.Printf("[Aiven/%s] Iniciando servicio: %s", c.project, service.Name)

		if err := c.client.StartService(c.project, service); err != nil {
			log.Printf(
				"[Aiven/%s] error iniciando servicio %s: %v",
				c.project, service.Name, err,
			)
			continue
		}
	}

	return nil
}
