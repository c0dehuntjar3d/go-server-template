package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var ErrDBNoURL = errors.New("database configuration: URL is required if other parameters are set")

const (
	defaultAppName    = "Go-Server"
	defaultAppVersion = "0.0.1"

	defaultHttpAddress     = "80"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second

	defaultLogLevel        = "info"
	defaultLogEncoding     = "console"
	defaultLogOutputPath   = "stderr"
	defaultLogErrorEnabled = true

	defaultDBConnectionTimeout  = 30
	defaultDBConnectionAttempts = 3
)

type (
	Config struct {
		App   *App
		Http  *HTTP
		Log   *Log
		DB    *DB
		Cache *Cache
	}

	App struct {
		Name    string
		Version string
	}

	HTTP struct {
		Address         string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
	}

	Log struct {
		Level        string
		Encoding     string
		OutputPath   string
		ErrorEnabled bool
	}

	DB struct {
		User               string
		Password           string
		Database           string
		Host               string
		SSL                string
		Port               int
		MaxConnection      int
		ConnectionTimeout  int
		ConnectionAttempts int
	}

	Cache struct {
		URL string
	}
)

var hdlOnce sync.Once

var config *Config

func LoadOrGetSingleton() (*Config, error) {
	hdlOnce.Do(func() {

		envPath := getEnvPath()
		err := godotenv.Load(*envPath)
		if err != nil {
			fmt.Println("Warning: .env file not found")
		}

		config = New()
	})

	return config, nil
}

func Default() *Config {
	return &Config{
		App: &App{
			Name:    defaultAppName,
			Version: defaultAppVersion,
		},
		Http: &HTTP{
			Address:         defaultHttpAddress,
			ReadTimeout:     defaultReadTimeout,
			WriteTimeout:    defaultWriteTimeout,
			ShutdownTimeout: defaultShutdownTimeout,
		},
		Log: &Log{
			Level:        defaultLogLevel,
			Encoding:     defaultLogEncoding,
			OutputPath:   defaultLogOutputPath,
			ErrorEnabled: defaultLogErrorEnabled,
		},
		DB:    &DB{},
		Cache: &Cache{},
	}
}

func New() *Config {
	return &Config{
		App: &App{
			Name:    getEnv("APP_NAME", defaultAppName),
			Version: getEnv("APP_VERSION", defaultAppVersion),
		},
		Http: &HTTP{
			Address:         getEnv("HTTP_ADDRESS", defaultHttpAddress),
			ReadTimeout:     getEnvAsDuration("HTTP_READ_TIMEOUT", defaultReadTimeout),
			WriteTimeout:    getEnvAsDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout),
			ShutdownTimeout: getEnvAsDuration("HTTP_SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
		},
		Log: &Log{
			Level:      getEnv("LOG_LEVEL", defaultLogLevel),
			Encoding:   getEnv("LOG_ENCODING", defaultLogEncoding),
			OutputPath: getEnv("LOG_OUTPUT_PATH", defaultLogOutputPath),
		},
		DB: &DB{
			User:               getEnv("DB_USER", ""),
			Password:           getEnv("DB_PASSWORD", ""),
			Database:           getEnv("DB_DATABASE", ""),
			SSL:                getEnv("DB_SSL", "disable"),
			Host:               getEnv("DB_HOST", "localhost"),
			MaxConnection:      getEnvAsInt("DB_MAX_CONNECTION", 0),
			Port:               getEnvAsInt("DB_PORT", 5432),
			ConnectionTimeout:  getEnvAsInt("DB_CONNECTION_TIMEOUT", defaultDBConnectionTimeout),
			ConnectionAttempts: getEnvAsInt("DB_CONNECTION_ATTEMPTS", defaultDBConnectionAttempts),
		},
		Cache: &Cache{
			URL: getEnv("REDIS_URL", ""),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	valStr := getEnv(name, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnvPath() *string {
	configPath := flag.String(
		"env",
		"./.env",
		"Path to configuration",
	)
	flag.Parse()

	return configPath
}
