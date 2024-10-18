ENV_FILE := .env 

include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))  

local-run: deps test build database migration run

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

# docker  

docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	@docker build -t $(DOCKER_IMAGE) . 

docker-run:
	@echo "Running Docker container $(DOCKER_CONTAINER)..."
	@docker run --env-file $(ENV_FILE) --name $(DOCKER_CONTAINER) -p $(HTTP_ADDRESS):$(HTTP_ADDRESS) $(DOCKER_IMAGE)

docker-clean:
	@echo "Stopping and removing Docker container $(DOCKER_CONTAINER)..."
	@docker stop $(DOCKER_CONTAINER) || true 
	@docker rm $(DOCKER_CONTAINER) || true

.PHONY: build test run clean deps
