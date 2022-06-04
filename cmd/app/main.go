package main

import (
	"context"
	"database/sql"
	"github.com/danielmichaels/storeman/internal/config"
	"github.com/danielmichaels/storeman/internal/server"
	"github.com/danielmichaels/storeman/internal/templates"
	"github.com/go-chi/httplog"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("failed to start storeman")
	}
}

// run is Storeman's entrypoint.
func run() error {
	cfg := config.AppConfig()
	logger := httplog.NewLogger("web-server", httplog.Options{
		JSON:     cfg.Logger.Json,
		Concise:  cfg.Logger.Concise,
		LogLevel: cfg.Logger.Level,
	})
	//db, err := openDB(cfg)
	//if err != nil {
	//	logger.Fatal().Err(err).Msg("failed to open database. exiting")
	//}
	logger.Info().Msg("database connection established")
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create a template cache")
	}

	app := server.Server{
		Config:   cfg,
		Logger:   logger,
		Template: templateCache,
	}

	err = app.Serve()
	if err != nil {
		app.Logger.Error().Err(err).Msg("server failed to start")
	}

	return nil
}

// openDB returns a sql connection pool
func openDB(cfg *config.Conf) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the
	// config struct
	db, err := sql.Open("sqlite3", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PingContext establishes a new connection to the database, passing in the
	// ctx as a parameter. If the connection couldn't be established within
	// 5 seconds, an error will be raised.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
