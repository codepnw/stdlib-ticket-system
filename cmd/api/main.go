package main

import (
	"log"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/pkg/database"
)

const envPath = ".env.example"

func main() {
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}
	
	db, err := database.ConnectPostgres(cfg.GetDBConnection())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	log.Println("Database Connected...")
}