package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/middleware"
	"github.com/codepnw/stdlib-ticket-system/internal/server"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
	jwttoken "github.com/codepnw/stdlib-ticket-system/pkg/jwt"
)

const envPath = ".env.example"

func main() {
	// Load Config
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	// Connect Database
	db, err := database.ConnectPostgres(cfg.GetDBConnection())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup Server
	serverCfg, err := setup(cfg, db)
	if err != nil {
		log.Fatal(err)
	}

	// Run Server
	if err := server.Run(serverCfg); err != nil {
		log.Fatal(err)
	}
}

func setup(cfg *config.EnvConfig, db *sql.DB) (*server.ServerConfig, error) {
	// Load Time Location
	location, _ := time.LoadLocation("Asia/Bangkok")

	// Database Transaction
	tx, err := database.NewTransaction(db)
	if err != nil {
		return nil, err
	}

	// Init JWT
	token, err := jwttoken.NewJWT(cfg.JWT.SecretKey, cfg.JWT.RefreshKey)
	if err != nil {
		return nil, err
	}

	// Mux Server
	mux := http.NewServeMux()

	// New Middleware
	mid := middleware.NewMiddleware(token)

	serverCfg := &server.ServerConfig{
		Location:   location,
		DB:         db,
		Tx:         tx,
		Mux:        mux,
		Addr:       ":8080",
		Token:      token,
		Middleware: mid,
	}
	return serverCfg, nil
}
