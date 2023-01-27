package main

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/mymorkkis/lets-go-further-json-api/internal/data"
	"github.com/mymorkkis/lets-go-further-json-api/internal/jsonlog"
	"github.com/mymorkkis/lets-go-further-json-api/internal/mailer"
)

const version = "1.0.0"

type application struct {
	version string
	config  *config
	logger  *jsonlog.Logger
	models  data.Models
	mailer  mailer.Mailer
	wg      sync.WaitGroup
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

	smtp := config.smtp

	app := &application{
		version: version,
		config:  config,
		logger:  logger,
		models:  data.NewModels(db),
		mailer:  mailer.New(smtp.host, smtp.port, smtp.username, smtp.password, smtp.sender),
	}

	err = app.serve()
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
