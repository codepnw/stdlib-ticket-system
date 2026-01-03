package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ConnectPostgres(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("db connect failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}
	return db, nil
}

// ==============  Transaction ===============

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) 
}

type txManager struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) (TxManager, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	return &txManager{db: db}, nil
}

func (t *txManager) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Fatalf("rollback: %v", err)
			}
		} else {
			cmErr := tx.Commit()
			err = cmErr
		}
	}()
	
	err = fn(tx)
	return err
}
