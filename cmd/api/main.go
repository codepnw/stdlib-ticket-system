package main

import (
	"log"
	"net/http"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
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

	// Database Transaction
	tx, err := database.NewTransaction(db)
	if err != nil {
		log.Fatal(err)
	}
	
	// Init JWT 
	token, err := jwttoken.NewJWT(cfg.JWT.SecretKey, cfg.JWT.RefreshKey)
	if err != nil {
		log.Fatal(err)
	}

	// Mux Server
	mux := http.NewServeMux()

	serverCfg := &server.ServerConfig{
		DB:    db,
		Tx:    tx,
		Mux:   mux,
		Addr:  ":8080",
		Token: token,
	}
	// Run Server
	if err := server.Run(serverCfg); err != nil {
		log.Fatal(err)
	}
}
