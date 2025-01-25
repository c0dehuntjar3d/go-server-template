ENV_FILE := .env

include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))  

startup: database deps test build migration run

env:
	@echo "Creating local env..."
	@cp env.example.env .env
	@cp env.example.env docker.env
	@echo ".env and docker.env files was created successfully!"
	@echo "Don't forget to update files with new variables!"

test: 
	@echo "Running tests..."
	@go test ./...

run:
	@echo "Running $(APP_NAME)..."
	@source $(ENV_FILE) && bin/$(APP_NAME)

build:
	@echo "Building $(APP_NAME)"
	@mkdir -p bin
	@go build -o bin/$(APP_NAME) cmd/server/main.go

deps:
	@echo "Installing dependencies..."
	@go mod tidy

database:
	@echo "Creating volume..."
	@docker volume create $(DB_VOLUME)
	@echo "Database loading..."
	@docker start $(APP_NAME)-db || docker run --name $(APP_NAME)-db -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_DATABASE) -v $(DB_VOLUME):/var/lib/postgresql/data -d -p $(DB_PORT):$(DB_PORT) postgres
	@echo "Database created"

migration:
	@echo "Applying migrations..."
	@migrate -path=migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_DATABASE)?sslmode=disable" -verbose up
	@echo "Migrations applied success!"

.PHONY: build test run clean deps
