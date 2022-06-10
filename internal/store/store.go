package store

import (
	"context"
	"database/sql"
	"github.com/danielmichaels/storeman/internal/config"
	"github.com/danielmichaels/storeman/internal/store/sqlite"
	"os"
	"path/filepath"
	"time"
)

type Store interface {
	ContainerInsert(title, notes string) (int, error)
	ContainerUpdate(title, notes string, id int) (int, error)
	ContainerGet(id int) (*sqlite.Container, error)
	ContainerGetAll() ([]*sqlite.Container, error)
	ContainerDelete(id int) error

	ItemInsert(fk int, name, description string, image []byte) (int, error)
	ItemGet(id int) (*sqlite.Item, error)
	ItemGetAllByContainer(id int) ([]*sqlite.Item, error)
	ItemDelete(id int) error
}

// OpenDB returns a sql connection pool
func OpenDB(cfg *config.Conf) (*sql.DB, error) {
	mkdir(filepath.Dir(cfg.Db.Dsn))
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

// mkdir creates a directory for the database if it does not already exist.
// mkdir is the functional equivalent of `mkdir -p` in POSIX compliant shell.
func mkdir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
