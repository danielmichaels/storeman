package config

import (
	"fmt"
	"github.com/joeshaw/envdecode"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("DB_NAME", "./data/store.db")

	cfg := AppConfig()

	if cfg.Db.Dsn != "./data/store.db" {
		t.Errorf("expected %q, got %q", "db", os.Getenv("POSTGRES_DB"))
	}
}
func ExampleAppConfig() {
	type exampleStruct struct {
		String string `env:"STRING"`
	}
	os.Setenv("STRING", "an example string!")

	var e exampleStruct
	err := envdecode.StrictDecode(&e)
	if err != nil {
		panic(err)
	}

	// if STRING is set, e.String will contain its value
	fmt.Println(e.String)

	// Output:
	// an example string!

}
