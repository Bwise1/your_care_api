package db

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
}

func New(dsn string) (*DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	return &DB{db}, nil
}

func (db *DB) RunInTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	// Begin a new transaction
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Ensure rollback if fn returns an error or panic occurs
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() // rollback on panic
			panic(p)          // re-throw panic after rollback
		} else if err != nil {
			_ = tx.Rollback() // rollback on error
		}
	}()

	// Run the provided function with the transaction
	if err = fn(tx); err != nil {
		return err // error will trigger rollback in defer
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
