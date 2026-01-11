package bookingusecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/codepnw/stdlib-ticket-system/internal/errs"
	"github.com/codepnw/stdlib-ticket-system/internal/features/booking"
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
			name:    "fail create items",
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
			// Setup
			uc, mockTx, mockBook, mockSeat := setup(t)

			// Mock FN
			tc.mockFn(mockTx, mockBook, mockSeat, tc.eventID, tc.seatIDs)

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

func TestGetBookingHistory(t *testing.T) {
	type testCase struct {
		name   string
		userID int64
		mockFn func(
			tx database.TxManager,
			mockBook bookingrepo.MockBookingRepository,
			mockSeat seatrepo.MockSeatRepository,
			userID int64,
		)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:   "success",
			userID: 1,
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64) {
				mockData := []booking.BookingHistoryResponse{
					{ID: "mock-uuid-1", EventDate: time.Now(), CreatedAt: time.Now()},
					{ID: "mock-uuid-2", EventDate: time.Now(), CreatedAt: time.Now()},
				}
				mockBook.EXPECT().GetHistory(gomock.Any(), userID).Return(mockData, nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:   "fail",
			userID: 1,
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64) {
				mockBook.EXPECT().GetHistory(gomock.Any(), userID).Return(nil, ErrMockDBError).Times(1)
			},
			expectedErr: ErrMockDBError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockTx, mockBook, mockSeat := setup(t)

			// Mock FN
			tc.mockFn(mockTx, mockBook, mockSeat, tc.userID)

			_, err := uc.GetBookingHistory(context.Background())

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCancelBooking(t *testing.T) {
	type testCase struct {
		name      string
		userID    int64
		bookingID string
		mockFn    func(
			tx database.TxManager,
			mockBook bookingrepo.MockBookingRepository,
			mockSeat seatrepo.MockSeatRepository,
			userID int64,
			bookingID string,
		)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:      "success",
			userID:    1,
			bookingID: "mock-uuid-1",
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64, bookingID string) {
				mockBookData := booking.Booking{ID: "mock-uuid-1", UserID: 1, Status: booking.StatusPending}
				mockBook.EXPECT().GetByID(gomock.Any(), bookingID).Return(mockBookData, nil).Times(1)

				mockBook.EXPECT().CancelBookingTx(gomock.Any(), gomock.Any(), mockBookData.ID).Return(nil).Times(1)

				mockSeat.EXPECT().CancelSeatsTx(gomock.Any(), gomock.Any(), mockBookData.ID).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:      "fail user other booking",
			userID:    2,
			bookingID: "mock-uuid-1",
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64, bookingID string) {
				mockBookData := booking.Booking{ID: "mock-uuid-1", UserID: 1, Status: booking.StatusCancelled}
				mockBook.EXPECT().GetByID(gomock.Any(), bookingID).Return(mockBookData, nil).Times(1)
			},
			expectedErr: errs.ErrCancelOtherBooking,
		},
		{
			name:      "fail cancel already",
			userID:    1,
			bookingID: "mock-uuid-1",
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64, bookingID string) {
				mockBookData := booking.Booking{ID: "mock-uuid-1", UserID: 1, Status: booking.StatusCancelled}
				mockBook.EXPECT().GetByID(gomock.Any(), bookingID).Return(mockBookData, nil).Times(1)
			},
			expectedErr: errs.ErrBookingIsCancel,
		},
		{
			name:      "fail cancel paid status",
			userID:    1,
			bookingID: "mock-uuid-1",
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64, bookingID string) {
				mockBookData := booking.Booking{ID: "mock-uuid-1", UserID: 1, Status: booking.StatusPaid}
				mockBook.EXPECT().GetByID(gomock.Any(), bookingID).Return(mockBookData, nil).Times(1)
			},
			expectedErr: errs.ErrBookingIsPaid,
		},
		{
			name:      "fail cancel booking",
			userID:    1,
			bookingID: "mock-uuid-1",
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64, bookingID string) {
				mockBookData := booking.Booking{ID: "mock-uuid-1", UserID: 1, Status: booking.StatusPending}
				mockBook.EXPECT().GetByID(gomock.Any(), bookingID).Return(mockBookData, nil).Times(1)

				mockBook.EXPECT().CancelBookingTx(gomock.Any(), gomock.Any(), mockBookData.ID).Return(ErrMockDBError).Times(1)
			},
			expectedErr: ErrMockDBError,
		},
		{
			name:      "fail cancel seats",
			userID:    1,
			bookingID: "mock-uuid-1",
			mockFn: func(tx database.TxManager, mockBook bookingrepo.MockBookingRepository, mockSeat seatrepo.MockSeatRepository, userID int64, bookingID string) {
				mockBookData := booking.Booking{ID: "mock-uuid-1", UserID: 1, Status: booking.StatusPending}
				mockBook.EXPECT().GetByID(gomock.Any(), bookingID).Return(mockBookData, nil).Times(1)

				mockBook.EXPECT().CancelBookingTx(gomock.Any(), gomock.Any(), mockBookData.ID).Return(nil).Times(1)
				
				mockSeat.EXPECT().CancelSeatsTx(gomock.Any(), gomock.Any(), mockBookData.ID).Return(ErrMockDBError).Times(1)
			},
			expectedErr: ErrMockDBError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uc, mockTx, mockBook, mockSeat := setup(t)

			// Mock FN
			tc.mockFn(mockTx, mockBook, mockSeat, tc.userID, tc.bookingID)

			err := uc.CancelBooking(context.Background(), tc.bookingID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func setup(t *testing.T) (bookingusecase.BookingUsecase, mockTx, bookingrepo.MockBookingRepository, seatrepo.MockSeatRepository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loc, _ := time.LoadLocation("Asia/Bangkok")

	mockTx := mockTx{}
	mockBook := bookingrepo.NewMockBookingRepository(ctrl)
	mockSeat := seatrepo.NewMockSeatRepository(ctrl)
	uc := bookingusecase.NewBookingUsecase(loc, mockTx, mockBook, mockSeat)

	return uc, mockTx, *mockBook, *mockSeat
}
