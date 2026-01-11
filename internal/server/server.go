package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	bookinghandler "github.com/codepnw/stdlib-ticket-system/internal/features/booking/handler"
	bookingrepo "github.com/codepnw/stdlib-ticket-system/internal/features/booking/repo"
	bookingusecase "github.com/codepnw/stdlib-ticket-system/internal/features/booking/usecase"
	eventhandler "github.com/codepnw/stdlib-ticket-system/internal/features/event/handler"
	eventrepo "github.com/codepnw/stdlib-ticket-system/internal/features/event/repo"
	eventusecase "github.com/codepnw/stdlib-ticket-system/internal/features/event/usecase"
	seatrepo "github.com/codepnw/stdlib-ticket-system/internal/features/seat/repo"
	userhandler "github.com/codepnw/stdlib-ticket-system/internal/features/user/handler"
	userrepo "github.com/codepnw/stdlib-ticket-system/internal/features/user/repo"
	userusecase "github.com/codepnw/stdlib-ticket-system/internal/features/user/usecase"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
	jwttoken "github.com/codepnw/stdlib-ticket-system/pkg/jwt"
)

type ServerConfig struct {
	Location *time.Location
	DB       *sql.DB
	Mux      *http.ServeMux
	Tx       database.TxManager
	Addr     string
	Token    jwttoken.JWTToken
}

func Run(cfg *ServerConfig) error {
	cfg.eventRoutes()
	cfg.userRoutes()
	cfg.bookingRoutes()

	log.Println("server running...")

	if err := http.ListenAndServe(cfg.Addr, cfg.Mux); err != nil {
		return err
	}
	return nil
}

func (cfg ServerConfig) eventRoutes() {
	seatRepo := seatrepo.NewSeatRepository(cfg.DB)
	eventRepo := eventrepo.NewEventRepository(cfg.DB)
	uc := eventusecase.NewEventUsecase(cfg.Tx, eventRepo, seatRepo)
	handler := eventhandler.NewEventHandler(uc)

	cfg.Mux.HandleFunc("POST /events", handler.CreateEvent)
	cfg.Mux.HandleFunc("GET /events", handler.GetAllEvents)
	cfg.Mux.HandleFunc("GET /events/{event_id}", handler.GetEventByID)
	cfg.Mux.HandleFunc("GET /events/{event_id}/seats", handler.GetSeatsByEventID)
}

func (cfg ServerConfig) userRoutes() {
	repo := userrepo.NewUserRepository(cfg.DB)
	uc := userusecase.NewUserUsecase(cfg.Tx, cfg.Token, repo)
	handler := userhandler.NewUserHandler(uc)

	cfg.Mux.HandleFunc("POST /register", handler.Register)
	cfg.Mux.HandleFunc("POST /login", handler.Login)
}

func (cfg ServerConfig) bookingRoutes() {
	bookRepo := bookingrepo.NewBookingRepository(cfg.DB)
	seatRepo := seatrepo.NewSeatRepository(cfg.DB)
	uc := bookingusecase.NewBookingUsecase(cfg.Location, cfg.Tx, bookRepo, seatRepo)
	handler := bookinghandler.NewBookingHandler(uc)

	cfg.Mux.HandleFunc("POST /bookings", handler.CreateBooking)
	cfg.Mux.HandleFunc("GET /bookings/me", handler.GetBookingHistory)
	cfg.Mux.HandleFunc("POST /bookings/cancel", handler.CancelBooking)
}
