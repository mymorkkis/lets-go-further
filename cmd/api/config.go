package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	port int
	env  string
}

func ensureRequiredConfigProvided() error {
	requiredEnvVars := [2]string{"API_PORT", "API_ENV"}

	missingVars := []string{}

	for _, env := range requiredEnvVars {
		_, ok := os.LookupEnv(env)
		if !ok {
			missingVars = append(missingVars, env)
		}
	}

	if len(missingVars) > 0 {
		return errors.New(
			fmt.Sprintf(
				"Required env vars not provided: %s",
				strings.Join(missingVars, ", "),
			),
		)
	}

	return nil
}

func NewConfig() (*Config, error) {
	if err := ensureRequiredConfigProvided(); err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(os.Getenv("API_PORT"), 10, 64)
	if err != nil {
		return nil, err
	}

	c := Config{
		env:  os.Getenv("API_ENV"),
		port: int(port),
	}

	return &c, nil
}
