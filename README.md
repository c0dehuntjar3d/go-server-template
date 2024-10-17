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
APP_NAME=template-server
APP_VERSION=0.0.1

##HTTP server settings
HTTP_ADDRESS=8080
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=10s
HTTP_SHUTDOWN_TIMEOUT=5s

##Logging settings
LOG_LEVEL=debug
LOG_ENCODING=json
LOG_OUTPUT_PATH=/var/log/app.log
LOG_ERROR_ENABLED=true

##Database settings
DB_URL=postgresql://user:password@localhost:5432/mydb?sslmode=disable
DB_CONNECTION_TIMEOUT=60
DB_CONNECTION_ATTEMPTS=5

##Cache settings
REDIS_URL=redis://<user>:<pass>@localhost:6379/<db>
```

### Commands

To run your application with a specific `.env` file, use the `-env` flag:

```bash
go run main.go -env config/.env
```

### License

This project is licensed under the MIT License. See the `LICENSE` file for details.
