package server

import (
	"database/sql"
	"net/http"

	eventhandler "github.com/codepnw/stdlib-ticket-system/internal/features/event/handler"
	eventrepo "github.com/codepnw/stdlib-ticket-system/internal/features/event/repo"
	eventusecase "github.com/codepnw/stdlib-ticket-system/internal/features/event/usecase"
	seatrepo "github.com/codepnw/stdlib-ticket-system/internal/features/seat/repo"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
)

type ServerConfig struct {
	DB   *sql.DB
	Mux  *http.ServeMux
	Tx   database.TxManager
	Addr string
}

func Run(cfg *ServerConfig) error {
	cfg.eventRoutes()

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
}
