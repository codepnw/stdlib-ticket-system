package bookingusecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	bookingrepo "github.com/codepnw/stdlib-ticket-system/internal/features/booking/repo"
	bookingusecase "github.com/codepnw/stdlib-ticket-system/internal/features/booking/usecase"
	"github.com/codepnw/stdlib-ticket-system/internal/features/seat"
	seatrepo "github.com/codepnw/stdlib-ticket-system/internal/features/seat/repo"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ErrMockDBError = errors.New("db error")

type mockTx struct{}

func (m mockTx) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	return fn(nil)
}

func TestCreateBooking(t *testing.T) {
	type testCase struct {
		name    string
		eventID int64
		seatIDs []int64
		mockFn  func(
			tx database.TxManager,
			mockBook bookingrepo.MockBookingRepository,
			mockSeat seatrepo.MockSeatRepository,
			eventID int64,
			seatIDs []int64,
		)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:    "success",
			eventID: 10,
			seatIDs: []int64{11, 20, 30},
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, eventID int64, seatIDs []int64) {
				mockSeats := []seat.Seat{
					{ID: 11, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 20, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 30, Price: 200, Status: seat.StatusAvailable},
				}
				mockSeat.EXPECT().GetSeatsForUpdateTx(gomock.Any(), gomock.Any(), seatIDs).Return(mockSeats, nil).Times(1)

				mockSeat.EXPECT().UpdateSeatsStatusTx(gomock.Any(), gomock.Any(), seatIDs, string(seat.StatusSold)).Return(nil).Times(1)

				mockBookID := "mock-uuid-1"
				mockBook.EXPECT().CreateBookingTx(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockBookID, nil).Times(1)

				mockBook.EXPECT().CreateBookingItemsTx(gomock.Any(), gomock.Any(), mockBookID, seatIDs).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:    "fail some seat not available",
			eventID: 10,
			seatIDs: []int64{11, 20, 30},
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, eventID int64, seatIDs []int64) {
				mockSeats := []seat.Seat{
					{ID: 11, Price: 100, Status: seat.StatusSold}, 
					{ID: 20, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 30, Price: 200, Status: seat.StatusAvailable},
				}
				mockSeat.EXPECT().GetSeatsForUpdateTx(gomock.Any(), gomock.Any(), seatIDs).Return(mockSeats, nil).Times(1)
			},
			expectedErr: errs.ErrSomeSeatNotAvailable,
		},
		{
			name:    "fail update seats",
			eventID: 10,
			seatIDs: []int64{11, 20, 30},
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, eventID int64, seatIDs []int64) {
				mockSeats := []seat.Seat{
					{ID: 11, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 20, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 30, Price: 200, Status: seat.StatusAvailable},
				}
				mockSeat.EXPECT().GetSeatsForUpdateTx(gomock.Any(), gomock.Any(), seatIDs).Return(mockSeats, nil).Times(1)

				mockSeat.EXPECT().UpdateSeatsStatusTx(gomock.Any(), gomock.Any(), seatIDs, string(seat.StatusSold)).Return(ErrMockDBError).Times(1)
			},
			expectedErr: ErrMockDBError,
		},
		{
			name:    "fail create booking",
			eventID: 10,
			seatIDs: []int64{11, 20, 30},
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, eventID int64, seatIDs []int64) {
				mockSeats := []seat.Seat{
					{ID: 11, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 20, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 30, Price: 200, Status: seat.StatusAvailable},
				}
				mockSeat.EXPECT().GetSeatsForUpdateTx(gomock.Any(), gomock.Any(), seatIDs).Return(mockSeats, nil).Times(1)

				mockSeat.EXPECT().UpdateSeatsStatusTx(gomock.Any(), gomock.Any(), seatIDs, string(seat.StatusSold)).Return(nil).Times(1)

				mockBook.EXPECT().CreateBookingTx(gomock.Any(), gomock.Any(), gomock.Any()).Return("", ErrMockDBError).Times(1)
			},
			expectedErr: ErrMockDBError,
		},
		{
			name:    "success",
			eventID: 10,
			seatIDs: []int64{11, 20, 30},
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, eventID int64, seatIDs []int64) {
				mockSeats := []seat.Seat{
					{ID: 11, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 20, Price: 100, Status: seat.StatusAvailable}, 
					{ID: 30, Price: 200, Status: seat.StatusAvailable},
				}
				mockSeat.EXPECT().GetSeatsForUpdateTx(gomock.Any(), gomock.Any(), seatIDs).Return(mockSeats, nil).Times(1)

				mockSeat.EXPECT().UpdateSeatsStatusTx(gomock.Any(), gomock.Any(), seatIDs, string(seat.StatusSold)).Return(nil).Times(1)

				mockBookID := "mock-uuid-1"
				mockBook.EXPECT().CreateBookingTx(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockBookID, nil).Times(1)

				mockBook.EXPECT().CreateBookingItemsTx(gomock.Any(), gomock.Any(), mockBookID, seatIDs).Return(ErrMockDBError).Times(1)
			},
			expectedErr: ErrMockDBError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock DI
			mockTx := mockTx{}
			mockBook := bookingrepo.NewMockBookingRepository(ctrl)
			mockSeat := seatrepo.NewMockSeatRepository(ctrl)
			// Mock FN
			tc.mockFn(mockTx, *mockBook, *mockSeat, tc.eventID, tc.seatIDs)

			// Usecase
			uc := bookingusecase.NewBookingUsecase(mockTx, mockBook, mockSeat)
			// Create Booking
			err := uc.CreateBooking(context.Background(), tc.eventID, tc.seatIDs)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
