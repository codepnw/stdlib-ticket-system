package errs

import "errors"

var (
	ErrEventNotFound         = errors.New("event not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrSeatNotFound          = errors.New("seat not found")
	ErrSomeSeatNotAvailable  = errors.New("some seats not available")
	
	// Bookings
	ErrBookingNotFound       = errors.New("booking not found")
	ErrCancelOtherBooking    = errors.New("you cannot cancel bookings")
	ErrBookingIsCancel       = errors.New("booking already cancelled")
	ErrBookingIsPaid         = errors.New("cannot cancel paid booking ")
)
