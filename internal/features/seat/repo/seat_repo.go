package seatrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/seat"
	"github.com/lib/pq"
)

//go:generate mockgen -source=seat_repo.go -destination=seat_repo_mock.go -package=seatrepo
type SeatRepository interface {
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]seat.Seat, error)

	// Transaction
	CreateSeatBatchTx(ctx context.Context, tx *sql.Tx, seats []seat.Seat) error
	GetSeatsForUpdateTx(ctx context.Context, tx *sql.Tx, seatIDs []int64) ([]seat.Seat, error)
	UpdateSeatsStatusTx(ctx context.Context, tx *sql.Tx, seatIDs []int64, status string) error
	CancelSeatsTx(ctx context.Context, tx *sql.Tx, bookingID string) error
}

type seatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) CreateSeatBatchTx(ctx context.Context, tx *sql.Tx, seats []seat.Seat) error {
	if len(seats) == 0 {
		return nil
	}

	valStrs := make([]string, 0, len(seats))
	valArgs := make([]any, 0, len(seats)*5)

	for i, seat := range seats {
		n := i * 5
		placeholders := fmt.Sprintf("($%d, $%d, $%d, $%d,$%d)", n+1, n+2, n+3, n+4, n+5)

		valStrs = append(valStrs, placeholders)
		valArgs = append(valArgs, seat.EventID, seat.SeatNumber, seat.Price, seat.Status, seat.Version)
	}

	query := "INSERT INTO seats (event_id, seat_number, price, status, version) VALUES %s"
	query = fmt.Sprintf(query, strings.Join(valStrs, ","))

	_, err := tx.ExecContext(ctx, query, valArgs...)
	if err != nil {
		return err
	}
	return nil
}

func (r *seatRepository) GetSeatsByEventID(ctx context.Context, eventID int64) ([]seat.Seat, error) {
	query := `
		SELECT id, event_id, seat_number, price, status, version
		FROM seats WHERE event_id = $1 ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []seat.Seat
	for rows.Next() {
		var s seat.Seat
		if err := rows.Scan(
			&s.ID,
			&s.EventID,
			&s.SeatNumber,
			&s.Price,
			&s.Status,
			&s.Version,
		); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *seatRepository) GetSeatsForUpdateTx(ctx context.Context, tx *sql.Tx, seatIDs []int64) ([]seat.Seat, error) {
	query := `SELECT id, status, price FROM seats WHERE id = ANY($1) FOR UPDATE`
	rows, err := tx.QueryContext(ctx, query, pq.Array(seatIDs))
	if err != nil {
		return nil, err
	}

	var seats []seat.Seat
	for rows.Next() {
		var s seat.Seat
		if err := rows.Scan(
			&s.ID,
			&s.Status,
			&s.Price,
		); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *seatRepository) UpdateSeatsStatusTx(ctx context.Context, tx *sql.Tx, seatIDs []int64, status string) error {
	query := `UPDATE seats SET status = $1 WHERE id = ANY($2)`
	res, err := tx.ExecContext(ctx, query, status, pq.Array(seatIDs))
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrSeatNotFound
	}
	return nil
}

func (r *seatRepository) CancelSeatsTx(ctx context.Context, tx *sql.Tx, bookingID string) error {
	query := `
		UPDATE seats SET status = 'AVAILABLE'
		WHERE id IN (
			SELECT seat_id FROM booking_items 
			WHERE booking_id = $1
		)
	`
	res, err := tx.ExecContext(ctx, query, bookingID)
	if err != nil {
		return err
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return errs.ErrBookingNotFound
	}
	return nil
}