package minecraft

import (
	"fmt"

	"github.com/iverly/go-mcping/mcping"
)

type Checker struct {
	host string
	port uint16
}

func NewChecker(host string, port uint16) *Checker {
	return &Checker{
		host: host,
		port: port,
	}
}

func (c *Checker) IsOnline() (bool, error) {

	pinger := mcping.NewPinger()

	_, err := pinger.Ping(c.host, c.port)
	if err != nil {
		return false, fmt.Errorf("error ejecutando ping al servidor: %w", err)
	}
    
	return true, nil
}