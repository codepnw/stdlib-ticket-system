package eventusecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/features/event"
	eventrepo "github.com/codepnw/stdlib-ticket-system/internal/features/event/repo"
	"github.com/codepnw/stdlib-ticket-system/internal/features/seat"
	seatrepo "github.com/codepnw/stdlib-ticket-system/internal/features/seat/repo"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, req event.CreateEventReq) error
	GetEventByID(ctx context.Context, eventID int64) (event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]seat.Seat, error)
}

type eventUsecase struct {
	tx        database.TxManager
	eventRepo eventrepo.EventRepository
	seatRepo  seatrepo.SeatRepository
}

func NewEventUsecase(tx database.TxManager, eventRepo eventrepo.EventRepository, seatRepo seatrepo.SeatRepository) EventUsecase {
	return &eventUsecase{
		tx:        tx,
		eventRepo: eventRepo,
		seatRepo:  seatRepo,
	}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, req event.CreateEventReq) error {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	err := u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		eventID, err := u.eventRepo.CreateEventTx(ctx, tx, event.Event{
			Name:      req.Name,
			EventDate: req.EventDate,
			IsActive:  req.IsActive,
		})
		if err != nil {
			return err
		}

		seats := make([]seat.Seat, 0)

		for _, zone := range req.Zones {
			for i := 1; i <= zone.SeatsPerRow; i++ {
				s := seat.Seat{
					EventID:    eventID,
					SeatNumber: fmt.Sprintf("%s%d", zone.ZoneName, i),
					Price:      zone.Price,
					Status:     seat.StatusAvailable,
					Version:    1,
				}
				seats = append(seats, s)
			}
		}

		if len(seats) > 0 {
			if err := u.seatRepo.CreateSeatBatchTx(ctx, tx, seats); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (u *eventUsecase) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	return u.eventRepo.GetAllEvents(ctx)
}

func (u *eventUsecase) GetEventByID(ctx context.Context, eventID int64) (event.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	return u.eventRepo.GetEventByID(ctx, eventID)
}

func (u *eventUsecase) GetSeatsByEventID(ctx context.Context, eventID int64) ([]seat.Seat, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	return u.seatRepo.GetSeatsByEventID(ctx, eventID)
}
