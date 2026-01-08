include .env.example

run:
	@go run cmd/api/main.go
	
test:
	@go test ./internal/features/booking/usecase -cover

docker-up:
	@docker compose --env-file=.env.example up -d

# --------------- MIGRATIONS -----------------
# 
MIGRATE_PART = "pkg/database/migrations"
MIGRATE_DB = "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_PART) -seq $(name)
	
migrate-up:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PART) up
	
migrate-down:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PART) down 1
	
migrate-force:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_PART) force $(version)
#	
# --------------- END MIGRATIONS -----------------

