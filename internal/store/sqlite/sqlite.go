package sqlite

import (
	"database/sql"
	"fmt"
)

type store struct {
	DB *sql.DB
}

func NewStore(DB *sql.DB) *store {
	return &store{DB: DB}
}

func (s store) Example() {
	fmt.Println("example")
}
