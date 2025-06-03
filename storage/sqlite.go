package storage

import (
	"context"
	"database/sql"
	"sync"

	"github.com/go-faster/errors"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gotd/td/session"
)

type SQLiteSessionStorage struct {
	DB   *sql.DB
	Path string
	mux  sync.Mutex
}

func NewSQLiteSessionStorage(path string) (*SQLiteSessionStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlite")
	}

	_, err = db.Exec(`PRAGMA journal_mode = WAL`)
	if err != nil {
		return nil, errors.Wrap(err, "set journal mode")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY CHECK (id = 0),
			data BLOB NOT NULL
		)
	`)
	if err != nil {
		return nil, errors.Wrap(err, "create table")
	}

	return &SQLiteSessionStorage{
		DB:   db,
		Path: path,
	}, nil
}

func (s *SQLiteSessionStorage) LoadSession(_ context.Context) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var data []byte
	err := s.DB.QueryRow(`SELECT data FROM sessions WHERE id = 0`).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, session.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "query row")
	}
	return data, nil
}

func (s *SQLiteSessionStorage) StoreSession(_ context.Context, data []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, err := s.DB.Exec(`
		INSERT INTO sessions (id, data) VALUES (0, ?) 
		ON CONFLICT(id) DO UPDATE SET data=excluded.data
	`, data)
	if err != nil {
		return errors.Wrap(err, "exec insert/update")
	}
	return nil
}
