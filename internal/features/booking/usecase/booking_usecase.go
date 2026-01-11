package bookingusecase

import (
	"context"
	"database/sql"
	"log"
	"time"

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
	GetBookingHistory(ctx context.Context) ([]displayBookingHistory, error)
	CancelBooking(ctx context.Context, bookingID string) error
}

type bookingUsecase struct {
	location *time.Location
	tx       database.TxManager
	bookRepo bookingrepo.BookingRepository
	seatRepo seatrepo.SeatRepository
}

func NewBookingUsecase(location *time.Location, tx database.TxManager, bookRepo bookingrepo.BookingRepository, seatRepo seatrepo.SeatRepository) BookingUsecase {
	return &bookingUsecase{
		location: location,
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

type displayBookingHistory struct {
	ID          string  `json:"id" `
	EventName   string  `json:"event_name" `
	TotalAmount float64 `json:"total_amount" `
	Status      string  `json:"status" `
	SeatNumbers string  `json:"seat_numbers"`
	EventDate   string  `json:"event_date" `
	CreatedAt   string  `json:"created_at" `
}

func (u *bookingUsecase) GetBookingHistory(ctx context.Context) ([]displayBookingHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	// TODO: Get UserID From Context Later
	userID := int64(1)

	history, err := u.bookRepo.GetHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []displayBookingHistory
	timeFormat := time.DateTime

	for _, h := range history {
		result = append(result, displayBookingHistory{
			ID:          h.ID,
			EventName:   h.EventName,
			TotalAmount: h.TotalAmount,
			Status:      h.Status,
			SeatNumbers: h.SeatNumbers,
			// Format time.Time -> Asia/Bangkok "2006-01-02 15:04:05"
			EventDate: h.EventDate.In(u.location).Format(timeFormat),
			CreatedAt: h.CreatedAt.In(u.location).Format(timeFormat),
		})
	}
	return result, nil
}

func (u *bookingUsecase) CancelBooking(ctx context.Context, bookingID string) error {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	// TODO: Get UserID From Context
	userID := int64(1)

	// 1. Check Booking ID
	bookData, err := u.bookRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}
	// 2. Check Ownership (userID != booking.UserID)
	if userID != bookData.UserID {
		return errs.ErrCancelOtherBooking
	}
	// 3. Check Cancelled Status
	if bookData.Status == booking.StatusCancelled {
		return errs.ErrBookingIsCancel
	}
	// 4. Check Paid Status cannot cancel
	if bookData.Status == booking.StatusPaid {
		return errs.ErrBookingIsPaid
	}

	return u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		// 5. Cancel Booking
		if err := u.bookRepo.CancelBookingTx(ctx, tx, bookData.ID); err != nil {
			return err
		}

		// 6. Cancel Seats
		if err := u.seatRepo.CancelSeatsTx(ctx, tx, bookData.ID); err != nil {
			return err
		}
		return nil
	})
}
