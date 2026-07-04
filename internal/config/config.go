package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Log      LogConfig      `mapstructure:"log"`
	Postgres PostgresConfig `mapstructure:"postgres"`
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
	// Output    string `mapstructure:"output" env:"LOG_OUTPUT"`
	// AddSource bool   `mapstructure:"add_source" env:"LOG_ADD_SOURCE"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"host" env:"POSTGRES_HOST"`
	Port            string        `mapstructure:"port" env:"POSTGRES_PORT"`
	User            string        `mapstructure:"user" env:"POSTGRES_USER"`
	Password        string        `mapstructure:"password" env:"POSTGRES_PASSWORD"`
	Db              string        `mapstructure:"db" env:"POSTGRES_DB"`
	SSLMode         string        `mapstructure:"sslmode" env:"POSTGRES_SSLMODE"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" env:"POSTGRES_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" env:"POSTGRES_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" env:"POSTGRES_CONN_MAX_LIFETIME"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=%d",
		url.QueryEscape(p.User),
		url.QueryEscape(p.Password),
		p.Host, p.Port,
		url.PathEscape(p.Db),
		p.SSLMode,
		int(p.ConnectTimeout.Seconds()),
	)
}

func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if c.HTTP.Addr == "" {
		return fmt.Errorf("http.addr is required")
	}
	if c.Postgres.Host == "" {
		return fmt.Errorf("postgres.host is required")
	}
	if c.Postgres.Port == "" {
		return fmt.Errorf("postgres.port is required")
	}
	if c.Postgres.Db == "" {
		return fmt.Errorf("postgres.db is required")
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

	if path != "" {
		v.SetConfigFile(path)
		ext := configFileExtension(path)
		if ext != "" {
			v.SetConfigType(ext)
		}
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("reading config: %w", err)
			}
		}
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("reading config: %w", err)
			}
		}

		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}
		v.SetConfigName("config." + env)
		if err := v.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("reading env config: %w", err)
			}
		}
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
		"postgres.host":                   "POSTGRES_HOST",
		"postgres.port":                   "POSTGRES_PORT",
		"postgres.user":                   "POSTGRES_USER",
		"postgres.password":               "POSTGRES_PASSWORD",
		"postgres.db":                     "POSTGRES_DB",
		"postgres.sslmode":                "POSTGRES_SSLMODE",
		"postgres.max_open_conns":         "POSTGRES_MAX_OPEN_CONNS",
		"postgres.max_idle_conns":         "POSTGRES_MAX_IDLE_CONNS",
		"postgres.conn_max_lifetime":      "POSTGRES_CONN_MAX_LIFETIME",
		"postgres.connect_timeout":        "POSTGRES_CONNECT_TIMEOUT",
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
		"postgres.host":                   "localhost",
		"postgres.port":                   "5432",
		"postgres.user":                   "taropood",
		"postgres.password":               "taropood",
		"postgres.db":                     "taropood",
		"postgres.sslmode":                "disable",
		"postgres.max_open_conns":         25,
		"postgres.max_idle_conns":         10,
		"postgres.conn_max_lifetime":      "5m",
		"postgres.connect_timeout":        "5s",
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