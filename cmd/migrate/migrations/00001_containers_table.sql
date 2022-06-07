-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS containers
(
    id         INTEGER NOT NULL primary key AUTOINCREMENT,
    created_at DATE DEFAULT (datetime('now', 'utc')),
    updated_at DATE DEFAULT (datetime('now', 'utc')),
    url        TEXT,
    title      TEXT,
    notes      TEXT,
    image      BLOB,
    qr_code    TEXT
);
CREATE TABLE IF NOT EXISTS items
(
    id          INTEGER NOT NULL primary key AUTOINCREMENT,
    name        TEXT,
    description TEXT,
    image       BLOB,
    created_at DATE DEFAULT (datetime('now', 'utc')),
    updated_at DATE DEFAULT (datetime('now', 'utc')),
    FOREIGN KEY (id) REFERENCES containers (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS containers;
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
