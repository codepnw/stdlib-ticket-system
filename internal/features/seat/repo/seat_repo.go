package seatrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/codepnw/stdlib-ticket-system/internal/features/seat"
)

type SeatRepository interface {
	CreateSeatBatchTx(ctx context.Context, tx *sql.Tx, seats []seat.Seat) error
}

type seatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (s *seatRepository) CreateSeatBatchTx(ctx context.Context, tx *sql.Tx, seats []seat.Seat) error {
	if len(seats) == 0 {
		return nil
	}

	valStrs := make([]string, 0, len(seats))
	valArgs := make([]any, 0, len(seats)*5)

	for i, seat := range seats {
		n := i * 5
		placeholders := fmt.Sprintf("($%d, $%d, $%d, $%d,$%d)", n+1, n+2, n+3, n+4, n+5)
		
		valStrs = append(valStrs, placeholders)
		valArgs = append(valArgs, seat.EventID, seat.SeatNumber, seat.Price, seat.Status, 1)
	}
	
	query := "INSERT INTO seats (event_id, seat_number, price, status, version) VALUES %s"
	query = fmt.Sprintf(query, strings.Join(valStrs, ","))
	
	_, err := tx.ExecContext(ctx, query, valArgs...)
	if err != nil {
		return err
	}
	return nil
}
