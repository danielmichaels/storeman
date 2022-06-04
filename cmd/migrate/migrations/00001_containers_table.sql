-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS containers
(
    id         INTEGER NOT NULL primary key AUTOINCREMENT,
    created_at DATE DEFAULT (datetime('now', 'utc')),
    updated_at DATE DEFAULT (datetime('now', 'utc')),
    url        TEXT,
    code       TEXT
);
CREATE TABLE IF NOT EXISTS items
(
    id          INTEGER NOT NULL,
    name        TEXT,
    description TEXT,
    FOREIGN KEY (id) REFERENCES containers (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS containers;
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
