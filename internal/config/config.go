package config

import (
	"github.com/joeshaw/envdecode"
	"log"
	"time"
)

type dbConf struct {
	Dsn string `env:"DSN,default=./data/store.db"`
}

type serverConf struct {
	Port         int           `env:"SERVER_PORT,default=9090"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,default=5s"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,default=10s"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,default=20s"`
}

type loggerConf struct {
	Level   string `env:"LOG_LEVEL,default=info"`
	Concise bool   `env:"LOG_CONCISE,default=true"`
	Json    bool   `env:"LOG_JSON,default=true"`
}

// Conf is the parent struct which holds the children configuration structs.
// Conf is used to keep config structs in isolated logical groups.
type Conf struct {
	Server serverConf
	Db     dbConf
	Logger loggerConf
}

// AppConfig runs the setup and install of the applications' configuration variables
// AppConfig is called wherever the config data is required.
func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
