package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/mymorkkis/lets-go-further-json-api/internal/data"
	"github.com/mymorkkis/lets-go-further-json-api/internal/jsonlog"
)

const version = "1.0.0"

type application struct {
	version string
	config  *config
	logger  *jsonlog.Logger
	models  data.Models
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	config, err := NewConfig()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(config)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	app := &application{
		version: version,
		config:  config,
		logger:  logger,
		models:  data.NewModels(db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  config.env,
	})
	err = server.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(config *config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.db.dsn)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(config.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.db.maxOpenConns)
	db.SetMaxIdleConns(config.db.maxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
