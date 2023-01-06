package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type db struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	port int
	env  string
	db   *db
}

func NewConfig() (*config, error) {
	if err := ensureRequiredConfigProvided(); err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(os.Getenv("API_PORT"), 10, 0)
	if err != nil {
		return nil, err
	}

	db, err := getDBConfig()
	if err != nil {
		return nil, err
	}

	c := config{
		env:  os.Getenv("API_ENV"),
		port: int(port),
		db:   db,
	}

	return &c, nil
}

func ensureRequiredConfigProvided() error {
	requiredEnvVars := [6]string{
		"API_PORT",
		"API_ENV",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_PORT",
		"POSTGRES_DB",
	}

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

func getDBConfig() (*db, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@json_api_db:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	maxOpenConns, err := getOptionalIntEnv("MAX_OPEN_CONNS", 50)
	if err != nil {
		return nil, err
	}

	maxIdleConns, err := getOptionalIntEnv("MAX_IDLE_CONNS", 50)
	if err != nil {
		return nil, err
	}

	maxIdleTime, err := getOptionalIntEnv("MAX_IDLE_TIME_MINS", 15)
	if err != nil {
		return nil, err
	}

	db := &db{
		dsn:          dsn,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		maxIdleTime:  fmt.Sprintf("%dm", maxIdleTime),
	}

	return db, nil
}

func getOptionalIntEnv(key string, defaultValue int) (int, error) {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue, nil
	}

	parsedInt, err := strconv.ParseInt(env, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(parsedInt), nil
}
