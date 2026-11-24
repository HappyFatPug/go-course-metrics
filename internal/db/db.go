package db

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

func (db *DB) Close(ctx context.Context) error {
	done := make(chan error, 1)

	go func() {
		done <- db.db.Close()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return errors.New("DB connection close timeout reached")
	}
}
