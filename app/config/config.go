package config

import (
	"errors"
	"flag"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoFilename         = errors.New("filename invalid")
	ErrCannotReadFilename = errors.New("cannot read from file")
	ErrDBNoURL            = errors.New("database configuration: URL is required if other parameters are set")
)

const (
	_defaultConfigPath            = "../config/application.yaml"
	_defaultConfigFlag            = "config"
	_defaultConfigFlagDescription = "Path to configuration"

	// HTTP
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
	_defaultAddress         = "80"

	// APPLICATION
	_defaultAppName = "Hellow world, Applicaiton!"
	_defaultVersion = "0.0.1"

	// LOGGER
	_defaultOutputPath      = "stderr" // output to console. output.log - output to file
	_defaultLogLevel        = "info"
	_defaultLogErrorEnabled = true
	_defaultEncoding        = "console" // otherwise json

	// DATABASE
	_defaultDBPoolMax            = 10
	_defaultDBConnectionTimeout  = 30
	_defaultDBConnectionAttempts = 3
)

type (
	Config struct {
		App  *App  `yaml:"app"`
		Http *HTTP `yaml:"http"`
		Log  *Log  `yaml:"logger"`
		DB   *DB   `yaml:"database"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Address         string        `yaml:"address" env:"HTTP_ADDRESS"`
		ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
		WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
	}

	Log struct {
		Level        string `yaml:"log_level"   env:"LOG_LEVEL"`
		Encoding     string `yaml:"encoding" env:"ENCODING"`
		OutputPath   string `yaml:"output_path" env:"OUTPUT_PATH"`
		ErrorEnabled bool   `yaml:"error_enabled"   env:"ERROR_ENABLED"`
	}

	DB struct {
		URL                string `yaml:"url"      env:"URL"`
		PoolMax            int    `yaml:"pool_max" env:"POOL_MAX"`
		ConnectionTimeout  int    `yaml:"connection_timeout" env:"CONNECTION_TIMEOUT"`
		ConnectionAttempts int    `yaml:"connection_attempts" env:"CONNECTION_ATTEMTS"`
	}
)

var hdlOnce sync.Once
var config *Config

func LoadOrGetSingleton() (*Config, error) {

	hdlOnce.Do(func() {
		configPath := getConfigPath()

		var err error
		config, err = New(*configPath)
		if err != nil {

			if errors.Is(err, ErrNoFilename) || errors.Is(err, ErrCannotReadFilename) {
				config = WithDefault()
			} else {
				panic(err)
			}
		}
	})

	return config, nil
}

func WithDefault() *Config {
	return &Config{
		App: &App{
			Name:    _defaultAppName,
			Version: _defaultVersion,
		},
		Http: &HTTP{
			ReadTimeout:     _defaultReadTimeout,
			WriteTimeout:    _defaultWriteTimeout,
			ShutdownTimeout: _defaultShutdownTimeout,
			Address:         _defaultAddress,
		},
		Log: &Log{
			Level:        _defaultLogLevel,
			ErrorEnabled: _defaultLogErrorEnabled,
			OutputPath:   _defaultOutputPath,
			Encoding:     _defaultEncoding,
		},
		DB: &DB{},
	}
}

func New(configPath string) (*Config, error) {
	if configPath == "" {
		return WithDefault(), ErrNoFilename
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return WithDefault(), ErrCannotReadFilename
	}

	var cfg *Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.mergeWithDefaultSettings()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func getConfigPath() *string {
	configPath := flag.String(
		_defaultConfigFlag,
		_defaultConfigPath,
		_defaultConfigFlagDescription,
	)
	flag.Parse()

	return configPath
}

func (c *Config) mergeWithDefaultSettings() error {
	if c.App != nil {
		c.App.mergeWithDefault()
	} else {
		c.App = &App{}
		c.App.mergeWithDefault()
	}

	if c.Http != nil {
		c.Http.mergeWithDefault()
	} else {
		c.Http = &HTTP{}
		c.Http.mergeWithDefault()
	}

	if c.Log != nil {
		c.Log.mergeWithDefault()
	} else {
		c.Log = &Log{}
		c.Log.mergeWithDefault()
	}

	if c.DB != nil {
		err := c.DB.mergeWithDefault()
		if err != nil {
			return err
		}
	} else {
		c.DB = &DB{}
		c.DB.mergeWithDefault()
	}

	return nil
}

func (db *DB) mergeWithDefault() error {
	if db.URL == "" && (db.PoolMax != 0 || db.ConnectionTimeout != 0 || db.ConnectionAttempts != 0) {
		return ErrDBNoURL
	}

	if db.URL != "" {
		if db.PoolMax == 0 {
			db.PoolMax = _defaultDBPoolMax
		}
		if db.ConnectionTimeout == 0 {
			db.ConnectionTimeout = _defaultDBConnectionTimeout
		}
		if db.ConnectionAttempts == 0 {
			db.ConnectionAttempts = _defaultDBConnectionAttempts
		}
	}

	return nil
}

func (app *App) mergeWithDefault() {
	if app.Name == "" {
		app.Name = _defaultAppName
	}
	if app.Version == "" {
		app.Version = _defaultVersion
	}
}

func (http *HTTP) mergeWithDefault() {
	if http.ReadTimeout == 0 {
		http.ReadTimeout = _defaultReadTimeout
	}
	if http.WriteTimeout == 0 {
		http.WriteTimeout = _defaultWriteTimeout
	}
	if http.ShutdownTimeout == 0 {
		http.ShutdownTimeout = _defaultShutdownTimeout
	}
	if http.Address == "" {
		http.Address = _defaultAddress
	}
}

func (log *Log) mergeWithDefault() {
	if log.Level == "" {
		log.Level = _defaultLogLevel
	}

	if log.Encoding == "" {
		log.Encoding = _defaultEncoding
	}

	if log.OutputPath == "" {
		log.OutputPath = _defaultOutputPath
	}
}
