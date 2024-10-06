package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaultConfig(t *testing.T) {
	defaultConfig := WithDefault()

	require.Equal(t, _defaultAppName, defaultConfig.App.Name)
	require.Equal(t, _defaultVersion, defaultConfig.App.Version)
	require.Equal(t, _defaultReadTimeout, defaultConfig.Http.ReadTimeout)
	require.Equal(t, _defaultWriteTimeout, defaultConfig.Http.WriteTimeout)
	require.Equal(t, _defaultShutdownTimeout, defaultConfig.Http.ShutdownTimeout)
	require.Equal(t, _defaultAddress, defaultConfig.Http.Address)
	require.Equal(t, _defaultLogLevel, defaultConfig.Log.Level)
	require.Equal(t, _defaultLogErrorEnabled, defaultConfig.Log.ErrorEnabled)
	require.Equal(t, 0, defaultConfig.DB.PoolMax)
	require.Equal(t, 0, defaultConfig.DB.ConnectionTimeout)
	require.Equal(t, 0, defaultConfig.DB.ConnectionAttempts)
}

func TestLoadConfigWithDefaults(t *testing.T) {
	yamlData := `
app:
  name: "MyApp"
http:
  read_timeout: 10s
logger:
  log_level: "DEBUG"
database:
  url: "postgres://user:pass@localhost:5432/dbname"
`

	tmpfile, err := os.CreateTemp("", "config.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(yamlData))
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	config, err := New(tmpfile.Name())
	require.NoError(t, err)

	require.Equal(t, "MyApp", config.App.Name)
	require.Equal(t, _defaultVersion, config.App.Version)
	require.Equal(t, 10*time.Second, config.Http.ReadTimeout)
	require.Equal(t, _defaultWriteTimeout, config.Http.WriteTimeout)
	require.Equal(t, _defaultShutdownTimeout, config.Http.ShutdownTimeout)
	require.Equal(t, _defaultAddress, config.Http.Address)
	require.Equal(t, "DEBUG", config.Log.Level)
	require.Equal(t, "postgres://user:pass@localhost:5432/dbname", config.DB.URL)
	require.Equal(t, _defaultDBPoolMax, config.DB.PoolMax)
	require.Equal(t, _defaultDBConnectionTimeout, config.DB.ConnectionTimeout)
	require.Equal(t, _defaultDBConnectionAttempts, config.DB.ConnectionAttempts)
}

func TestLoadConfigWithMissingDBURL(t *testing.T) {
	yamlData := `
database:
  pool_max: 5
  connection_timeout: 10
`

	tmpfile, err := os.CreateTemp("", "config.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(yamlData))
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	config, err := New(tmpfile.Name())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database configuration: URL is required if other parameters are set")
	require.Equal(t, 5, config.DB.PoolMax)
}

func TestLoadConfigWithAllFields(t *testing.T) {
	yamlData := `
app:
  name: "CompleteApp"
  version: "1.0.0"
http:
  read_timeout: 15s
  write_timeout: 5s
  shutdown_timeout: 4s
  address: ":8080"
logger:
  log_level: "INFO"
  error_enabled: false
database:
  url: "postgres://user:pass@localhost:5432/dbname"
  pool_max: 20
  connection_timeout: 30
  connection_attempts: 5
`

	tmpfile, err := os.CreateTemp("", "config.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(yamlData))
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	config, err := New(tmpfile.Name())
	require.NoError(t, err)

	require.Equal(t, "CompleteApp", config.App.Name)
	require.Equal(t, "1.0.0", config.App.Version)
	require.Equal(t, 15*time.Second, config.Http.ReadTimeout)
	require.Equal(t, 5*time.Second, config.Http.WriteTimeout)
	require.Equal(t, 4*time.Second, config.Http.ShutdownTimeout)
	require.Equal(t, ":8080", config.Http.Address)
	require.Equal(t, "INFO", config.Log.Level)
	require.Equal(t, false, config.Log.ErrorEnabled)
	require.Equal(t, "postgres://user:pass@localhost:5432/dbname", config.DB.URL)
	require.Equal(t, 20, config.DB.PoolMax)
	require.Equal(t, 30, config.DB.ConnectionTimeout)
	require.Equal(t, 5, config.DB.ConnectionAttempts)
}
