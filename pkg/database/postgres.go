package database

import (
	"context"
	"database/sql"
	"fmt"
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
