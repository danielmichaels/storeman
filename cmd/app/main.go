package main

import (
	"github.com/danielmichaels/storeman/internal/config"
	"github.com/danielmichaels/storeman/internal/server"
	"github.com/danielmichaels/storeman/internal/store"
	"github.com/danielmichaels/storeman/internal/store/sqlite"
	"github.com/danielmichaels/storeman/internal/templates"
	"github.com/go-chi/httplog"
	"github.com/go-playground/form/v4"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("failed to start storeman")
	}
}

// run is Storeman's entrypoint.
func run() error {
	cfg := config.AppConfig()

	logger := httplog.NewLogger("storeman-server", httplog.Options{
		JSON:     cfg.Logger.Json,
		Concise:  cfg.Logger.Concise,
		LogLevel: cfg.Logger.Level,
	})

	db, err := store.OpenDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open database. exiting")
	}
	logger.Info().Msg("database connection established")

	formDecoder := form.NewDecoder()

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create a template cache")
	}

	app := server.Server{
		Config:   cfg,
		Logger:   logger,
		Template: templateCache,
		Store:    sqlite.NewStore(db),
		Form:     formDecoder,
	}

	err = app.Serve()
	if err != nil {
		app.Logger.Error().Err(err).Msg("server failed to start")
	}

	return nil
}
