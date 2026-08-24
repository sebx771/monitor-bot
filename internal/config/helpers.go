package config

import (
	"os"
	"strconv"
)

func GetEnabledVar(name string) (bool, error) {

	variable,err:= strconv.ParseBool(os.Getenv(name))
	if err != nil {
		return false , err
	}

	return variable , nil
		
}