ENV_FILE := .env

include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))  

startup: deps test build database migration run

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
	@go build -o bin/$(APP_NAME) cmd/main.go

deps:
	@echo "Installing dependencies..."
	@go mod tidy

database:
	@echo "Creating volume..."
	@docker volume create $(POSTGRES_VOLUME)
	@echo "Database loading..."
	@docker start $(APP_NAME)-db || docker run --name $(APP_NAME)-db -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) -v $(POSTGRES_VOLUME):/var/lib/postgresql/data -d -p 5432:5432 postgres
	@echo "Database created"

migration:
	@echo "Applying migrations..."
	@migrate -path=migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" -verbose up
	@echo "Migrations applied success!"

.PHONY: build test run clean deps
