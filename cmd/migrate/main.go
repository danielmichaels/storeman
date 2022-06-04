package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/danielmichaels/storeman/internal/config"
	"github.com/danielmichaels/storeman/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	dialect              = "sqlite3"
	migrationsEmbed      = "migrations"
	migrationsFileSystem = "./cmd/migrate/migrations"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
)

func main() {
	cfg := config.AppConfig()

	if !cfg.Logger.Json {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	flags.Usage = usage
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
	}

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	goose.SetBaseFS(migrationFs)

	// Both `create` and `fix` are still conducted on the filesystem. All other migrations
	// are run from the embedded migration files.
	// see: https://github.com/pressly/goose#embedded-sql-migrations
	command := args[0]
	switch command {
	case "create":
		if err := goose.Run("create", nil, migrationsFileSystem, args[1:]...); err != nil {
			log.Error().Err(err).Msg("migrate create error")
		}
		// fix the migrations to use incremental system rather than isosec.
		if err := goose.Run("fix", nil, migrationsFileSystem, args[1:]...); err != nil {
			log.Error().Err(err).Msg("migrate fix error")
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, migrationsFileSystem); err != nil {
			log.Error().Err(err).Msg("migrate fix error")
		}
		return
	}

	db, err := store.OpenDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database. exiting")
	}

	defer db.Close()

	if err := goose.SetDialect(dialect); err != nil {
		log.Error().Err(err).Msg("migrate dialect error")
	}

	if err := goose.Run(command, db, migrationsEmbed, args[1:]...); err != nil {
		log.Error().Err(err).Msgf("migrate %q error", command)
	}
}

func usage() {
	fmt.Print(usagePrefix)
	flags.PrintDefaults()
	fmt.Print(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)
