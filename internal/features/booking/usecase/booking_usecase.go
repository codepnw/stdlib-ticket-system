package bookingusecase

import (
	"context"
	"database/sql"
	"log"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/booking"
	bookingrepo "github.com/codepnw/stdlib-ticket-system/internal/features/booking/repo"
	"github.com/codepnw/stdlib-ticket-system/internal/features/seat"
	seatrepo "github.com/codepnw/stdlib-ticket-system/internal/features/seat/repo"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
)

type BookingUsecase interface {
	CreateBooking(ctx context.Context, eventID int64, seatIDs []int64) error
}

type bookingUsecase struct {
	tx       database.TxManager
	bookRepo bookingrepo.BookingRepository
	seatRepo seatrepo.SeatRepository
}

func NewBookingUsecase(tx database.TxManager, bookRepo bookingrepo.BookingRepository, seatRepo seatrepo.SeatRepository) BookingUsecase {
	return &bookingUsecase{
		tx:       tx,
		bookRepo: bookRepo,
		seatRepo: seatRepo,
	}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, eventID int64, seatIDs []int64) error {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	return u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		// Get Seats
		seats, err := u.seatRepo.GetSeatsForUpdateTx(ctx, tx, seatIDs)
		if err != nil {
			log.Printf("get seats failed: %v", err)
			return err
		}
		// Validate Seats Len
		if len(seats) != len(seatIDs) {
			return errs.ErrSomeSeatNotAvailable
		}

		var totalAmount float64
		for _, s := range seats {
			if s.Status != seat.StatusAvailable {
				return errs.ErrSomeSeatNotAvailable
			}
			totalAmount += s.Price
		}

		// Update Seats Status
		if err := u.seatRepo.UpdateSeatsStatusTx(ctx, tx, seatIDs, string(seat.StatusSold)); err != nil {
			log.Printf("update seats failed: %v", err)
			return err
		}

		// Create Booking
		bookingID, err := u.bookRepo.CreateBookingTx(ctx, tx, booking.Booking{
			UserID:      1, // TODO: Get From Context Later
			EventID:     eventID,
			TotalAmount: totalAmount,
			Status:      booking.StatusPending,
		})
		if err != nil {
			log.Printf("create booking failed: %v", err)
			return err
		}

		// Create Booking Items
		if err := u.bookRepo.CreateBookingItemsTx(ctx, tx, bookingID, seatIDs); err != nil {
			log.Printf("create booking items failed: %v", err)
			return err
		}
		return nil
	})
}
