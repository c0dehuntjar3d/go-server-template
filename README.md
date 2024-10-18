# Golang template server with database and cache

## Clean Architecture in This Project

This project follows the principles of Clean Architecture, ensuring a clear separation of concerns between different layers of the application.

## Configuration

This project provides a configuration module for a Go application, designed to load configuration parameters from environment variables or a `.env` file using the `godotenv` package. The configuration covers various aspects, including application metadata, HTTP server settings, logging, database connection, and caching.

### Usage

#### Loading the Configuration

To load the configuration, call `LoadOrGetSingleton()` which will:

- Parse environment variables using the `godotenv` package.
- Apply default values if the environment variables are not set.
- Use a `.env` file, which can be specified using the `-env` flag or defaults to `./.env`.

```go
config, err := config.LoadOrGetSingleton()
if err != nil {
    log.Fatal(err)
}
```

#### Default Configuration

You can also load the default configuration without relying on environment variables:

```go
defaultConfig := config.Default()
```

#### Example `.env` File

Create an `.env` file with the following structure to customize the configuration for your environment:

```env
##Application settings
APP_VERSION=0.1.0-dev
APP_NAME=server

##HTTP settings
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=5s
HTTP_SHUTDOWN_TIMEOUT=3s
HTTP_ADDRESS=:3000

##LOG settings
LOG_LEVEL=debug
LOG_OUTPUT_PATH=stderr
LOG_ENCODING=console

##DATABASE settings
DB_USER=postgres
DB_PASSWORD=1234
DB_DATABASE=appdb
DB_MAX_CONS=10 
DB_PORT=5432
DB_SSL=disable
DB_HOST=localhost
DB_VOLUME=app-volume

##DOCKER settings
DOCKER_IMAGE=server_image
DOCKER_CONTAINER=server_container 
```

### Commands

To run application:

#### With Makefile

```bash
make startup
```

#### With Go built-in

```bash
go run cmd/main.go
```

### License

This project is licensed under the MIT License. See the `LICENSE` file for details.
