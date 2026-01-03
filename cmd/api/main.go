package main

import (
	"log"
	"net/http"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/server"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
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

	// Database Transaction
	tx, err := database.NewTransaction(db)
	if err != nil {
		log.Fatal(err)
	}

	// Mux Server
	mux := http.NewServeMux()

	serverCfg := &server.ServerConfig{
		DB:   db,
		Tx:   tx,
		Mux:  mux,
		Addr: ":8080",
	}
	// Run Server
	if err := server.Run(serverCfg); err != nil {
		log.Fatal(err)
	}
}
