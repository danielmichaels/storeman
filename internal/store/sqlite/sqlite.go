package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const (
	QueryCtx = 3 * time.Second
)

var (
	ErrNoRecord = errors.New("sqlite: no matching record found")
)

type Container struct {
	ID    int
	Title string
	Notes string
	// todo
	Location string
	// todo
	Image     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	// todo
	Url string
	// todo
	QRCode string
}

type Item struct {
	ID          int
	ContainerId uint64
	Name        string
	Description string
	// todo
	Image     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

type store struct {
	DB *sql.DB
}

func NewStore(DB *sql.DB) *store {
	return &store{DB: DB}
}

func (s store) ContainerInsert(title, notes string) (int, error) {
	stmt := `
		INSERT INTO containers (title, notes, url, qr_code)
		VALUES ($1, $2, $3, $4)
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, stmt, title, notes, "url", "qr_code")
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
func (s store) ContainerUpdate(title, notes string, id int) (int, error) {
	stmt := `
	UPDATE containers
	SET
		title = $1,
		notes = $2
	WHERE
		id = $3
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()
	var c Container
	args := []interface{}{title, notes, id}
	err := s.DB.QueryRowContext(ctx, stmt, args...).Scan(&c.ID)
	if err != nil {
		return 0, err
	}

	if c.ID == 0 {
		return 0, ErrNoRecord
	}
	return c.ID, nil
}
func (s store) ContainerGet(id int) (*Container, error) {
	stmt := `
SELECT id, title, notes, url, qr_code, created_at, updated_at
FROM containers
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, stmt, id)

	c := &Container{}

	err := row.Scan(&c.ID, &c.Title, &c.Notes, &c.Url, &c.QRCode, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}
	return c, nil
}
func (s store) ContainerGetAll() ([]*Container, error) {
	stmt := `
SELECT id, title, url, notes, created_at, updated_at
FROM containers
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	// todo create filters and replace magic numbers
	args := []any{"10", "0"}

	rows, err := s.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []*Container

	for rows.Next() {
		var container Container
		err := rows.Scan(
			&container.ID,
			&container.Title,
			&container.Url,
			&container.Notes,
			&container.CreatedAt,
			&container.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		containers = append(containers, &container)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	if containers == nil {
		containers = []*Container{}
	}
	return containers, nil
}

func (s store) ContainerDelete(id int) error {
	stmt := `
DELETE FROM 
	containers
WHERE 
	id = $1
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
func (s store) ItemInsert(fk int, name, description string, image []byte) (int, error) {
	stmt := `
		INSERT INTO items (container_id, name, description, image)
		VALUES ($1, $2, $3, $4)
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, stmt, fk, name, description, image)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s store) ItemUpdate(id int, name, description string, image []byte) (int, error) {
	stmt := `
	UPDATE items
	SET
		name = $2,
		description = $3,
		image = $4
	WHERE
		id = $1
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()
	var c Container
	args := []interface{}{id, name, description, image}
	err := s.DB.QueryRowContext(ctx, stmt, args...).Scan(&c.ID)
	if err != nil {
		return 0, err
	}

	if c.ID == 0 {
		return 0, ErrNoRecord
	}
	return c.ID, nil
}

func (s store) ItemGet(id int) (*Item, error) {

	stmt := `
SELECT id, name, description, image, created_at, updated_at
FROM items
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, stmt, id)

	i := &Item{}

	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.Image, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}
	return i, nil
}
func (s store) ItemGetAllByContainer(id int) ([]*Item, error) {
	stmt := `
SELECT i.id, i.name, i.description, i.image, i.created_at, i.updated_at, i.container_id
FROM items i
WHERE i.container_id = $1
ORDER BY i.updated_at DESC
LIMIT $2 OFFSET $3
`
	ctx, cancel := context.WithTimeout(context.Background(), QueryCtx)
	defer cancel()

	// todo create filters and replace magic numbers
	args := []any{id, "10", "0"}

	rows, err := s.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item

	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Image,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ContainerId,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	if items == nil {
		items = []*Item{}
	}
	return items, nil
}
func (s store) ItemDelete(id int) error {
	return nil
}
