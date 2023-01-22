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

type limiter struct {
	rps     float64
	burst   int
	enabled bool
}

type config struct {
	port    int
	env     string
	db      *db
	limiter *limiter
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

	limiter, err := getLimiterConfig()
	if err != nil {
		return nil, err
	}

	c := config{
		env:     os.Getenv("API_ENV"),
		port:    int(port),
		db:      db,
		limiter: limiter,
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

func getLimiterConfig() (*limiter, error) {
	rps, err := getOptionalFloat64Env("LIMITER_RPS", 2.0)
	if err != nil {
		return nil, err
	}

	burst, err := getOptionalIntEnv("LIMITER_BURST", 4)
	if err != nil {
		return nil, err
	}

	limiterEnabled := true
	enabled := strings.ToLower(os.Getenv("LIMITER_ENABLED"))
	if enabled == "false" || enabled == "f" {
		limiterEnabled = false
	}

	limiter := &limiter{
		rps:     rps,
		burst:   burst,
		enabled: limiterEnabled,
	}

	return limiter, nil
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

func getOptionalFloat64Env(key string, defaultValue float64) (float64, error) {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue, nil
	}

	parsedFloat, err := strconv.ParseFloat(key, 64)
	if err != nil {
		return 0.0, err
	}

	return parsedFloat, nil
}
