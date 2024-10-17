ENV_FILE := .env 

include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))  

full-run: deps test build run

test: 
	@echo "Running tests..."
	@go test ./...

run:
	@echo "Running $(APP_NAME)..."
	@source $(ENV_FILE) && bin/${APP_NAME}

build:
	@echo "Building ${APP_NAME}"
	@mkdir -p bin
	@go build -o bin/${APP_NAME} cmd/main.go

deps:
	@echo "Installing dependencies..."
	@go mod tidy

docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	@docker build -t $(DOCKER_IMAGE) . 

docker-run:
	@echo "Running Docker container $(DOCKER_CONTAINER)..."
	@docker run --env-file $(ENV_FILE) --name $(DOCKER_CONTAINER) -p ${HTTP_ADDRESS}:${HTTP_ADDRESS} $(DOCKER_IMAGE)

docker-clean:
	@echo "Stopping and removing Docker container $(DOCKER_CONTAINER)..."
	@docker stop $(DOCKER_CONTAINER) || true 
	@docker rm $(DOCKER_CONTAINER) || true


.PHONY: build test run clean deps
