package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	HTTP       HTTPConfig       `mapstructure:"http"`
	Log        LogConfig        `mapstructure:"log"`
}

type AppConfig struct {
	Name    string `mapstructure:"name" env:"APP_NAME"`
	Env     string `mapstructure:"env" env:"APP_ENV"`
	Version string `mapstructure:"version"`
}

type HTTPConfig struct {
	Addr            string        `mapstructure:"addr" env:"HTTP_ADDR"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" env:"HTTP_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
}

type LogConfig struct {
	Level     string `mapstructure:"level" env:"LOG_LEVEL"`
	Format    string `mapstructure:"format" env:"LOG_FORMAT"`
	Output    string `mapstructure:"output" env:"LOG_OUTPUT"`
	AddSource bool   `mapstructure:"add_source" env:"LOG_ADD_SOURCE"`
}

// Validate checks required configuration and returns error on first failure.
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.HTTP.Addr == "" {
		return fmt.Errorf("http.addr is required")
	}

	return nil
}

// LoadConfig loads configuration from YAML files and environment variables.
// Priority (lowest to highest):
// 1. Hardcoded defaults
// 2. configs/config.yaml
// 3. configs/config.{env}.yaml (if APP_ENV is set)
// 4. Environment variables
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// 12-Factor III: Env vars take highest priority.
	// Viper maps nested keys to env vars using KEY_DELIMITER _ replacer.
	// e.g., "http.addr" → env "HTTP_ADDR"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 12-Factor: All config values must be overridable from env.
	// Bind explicit env vars for clarity and documentation.
	bindEnvKeys(v)

	// Set hardcoded defaults (lowest priority)
	setDefaults(v)

	// Load config file(s)
	if path != "" {
		// Explicit config file path
		v.SetConfigFile(path)
		ext := configFileExtension(path)
		if ext != "" {
			v.SetConfigType(ext)
		}
	} else {
		// Default search paths
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")

		// Load environment-specific overlay: config.dev.yaml
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}
		v.SetConfigName("config." + env)
		// Merge if exists (ignores file not found)
		_ = v.MergeInConfig()
	}

	// Read main config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
		// Config file is optional if all values come from env vars
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func bindEnvKeys(v *viper.Viper) {
	// Explicitly bind env vars to config keys for clarity.
	// Each binding: v.BindEnv("config.key", "ENV_VAR_NAME")
	envBindings := map[string]string{
		"app.name":                        "APP_NAME",
		"app.env":                         "APP_ENV",
		"http.addr":                       "HTTP_ADDR",
		"http.read_timeout":               "HTTP_READ_TIMEOUT",
		"http.write_timeout":              "HTTP_WRITE_TIMEOUT",
		"http.shutdown_timeout":           "HTTP_SHUTDOWN_TIMEOUT",
		"log.level":                       "LOG_LEVEL",
		"log.format":                      "LOG_FORMAT",
		"log.output":                      "LOG_OUTPUT",
		"log.add_source":                  "LOG_ADD_SOURCE",
	}

	for key, env := range envBindings {
		_ = v.BindEnv(key, env)
	}
}

func setDefaults(v *viper.Viper) {
	defaults := map[string]interface{}{
		"app.name":                        "taropood",
		"app.env":                         "dev",
		"http.addr":                       ":8091",
		"http.read_timeout":               "10s",
		"http.write_timeout":              "10s",
		"http.shutdown_timeout":           "15s",
		"log.level":                       "info",
		"log.format":                      "json",
		"log.output":                      "stdout",
		"log.add_source":                  false,
	
	}

	for key, val := range defaults {
		v.SetDefault(key, val)
	}
}

func configFileExtension(path string) string {
	if idx := strings.LastIndex(path, "."); idx >= 0 {
		return path[idx+1:]
	}
	return ""
}