package config

import (
	"errors"
	"os"
)

// Config holds all environment-driven configuration for control-api.
type Config struct {
	Port          string
	DBUrl         string
	RedisURL      string
	OtelEndpoint  string
	JWTSecret     string
}

// Load reads configuration from environment variables.
// Returns an error if any required field is missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnvOrDefault("CONTROL_API_PORT", "8080"),
		DBUrl:        os.Getenv("SUPABASE_DB_URL"),
		RedisURL:     os.Getenv("REDIS_URL"),
		OtelEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}

	var missing []string
	if cfg.DBUrl == "" {
		missing = append(missing, "SUPABASE_DB_URL")
	}
	if cfg.RedisURL == "" {
		missing = append(missing, "REDIS_URL")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return nil, errors.New("missing required env vars: " + joinStrings(missing))
	}

	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func joinStrings(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
