package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetEnabledVar(name string) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return false, nil
	}

	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(
			"la variable %s debe ser un booleano (true/false), se recibió %q",
			name,
			value,
		)
	}

	return enabled, nil
}
