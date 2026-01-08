package bookingrepo

import (
	"context"
	"database/sql"

	"github.com/codepnw/stdlib-ticket-system/internal/features/booking"
	"github.com/lib/pq"
)

type BookingRepository interface {
	CreateBookingTx(ctx context.Context, tx *sql.Tx, input booking.Booking) (string, error)
	CreateBookingItemsTx(ctx context.Context, tx *sql.Tx, bookingID string, seatIDs []int64) error
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBookingTx(ctx context.Context, tx *sql.Tx, input booking.Booking) (string, error) {
	query := `
		INSERT INTO bookings (user_id, event_id, total_amount, status)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	var id string

	err := tx.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.EventID,
		input.TotalAmount,
		input.Status,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *bookingRepository) CreateBookingItemsTx(ctx context.Context, tx *sql.Tx, bookingID string, seatIDs []int64) error {
	query := `
		INSERT INTO booking_items (booking_id, seat_id)
		SELECT $1, UNNEST($2::BIGINT[])
	`
	_, err := tx.ExecContext(ctx, query, bookingID, pq.Array(seatIDs))
	if err != nil {
		return err
	}
	return nil
}
