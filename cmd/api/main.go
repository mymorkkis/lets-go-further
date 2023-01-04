package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type application struct {
	version string
	config  *Config
	logger  *log.Logger
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	config, err := NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		version: version,
		config:  config,
		logger:  logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", config.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}
