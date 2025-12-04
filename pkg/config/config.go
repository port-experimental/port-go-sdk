package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config captures the high-level SDK settings. Additional per-client
// options will layer on top of this struct.
type Config struct {
	Region       string // eu (default) or us
	BaseURL      string // optional manual override
	ClientID     string
	ClientSecret string
	APIToken     string
}

// Load builds a Config from environment variables. If dotenvPath is non-empty,
// the file will be loaded before reading process environment variables.
func Load(dotenvPath string) (Config, error) {
	if dotenvPath != "" {
		_ = godotenv.Load(dotenvPath)
	} else {
		_ = godotenv.Load()
	}

	cfg := Config{
		Region:       envOr("PORT_REGION", "eu"),
		BaseURL:      os.Getenv("PORT_BASE_URL"),
		ClientID:     strings.TrimSpace(os.Getenv("PORT_CLIENT_ID")),
		ClientSecret: strings.TrimSpace(os.Getenv("PORT_CLIENT_SECRET")),
		APIToken:     strings.TrimSpace(os.Getenv("PORT_ACCESS_TOKEN")),
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate ensures required auth values are present.
func (c Config) Validate() error {
	if c.APIToken != "" {
		return nil
	}
	if c.ClientID == "" || c.ClientSecret == "" {
		return errors.New("set PORT_ACCESS_TOKEN or PORT_CLIENT_ID/PORT_CLIENT_SECRET")
	}
	return nil
}

// BaseEndpoint resolves the API base URL considering region + overrides.
func (c Config) BaseEndpoint() string {
	if strings.TrimSpace(c.BaseURL) != "" {
		return strings.TrimRight(c.BaseURL, "/")
	}
	switch strings.ToLower(c.Region) {
	case "us":
		return "https://api.us.port.io"
	case "eu", "":
		return "https://api.port.io"
	default:
		return fmt.Sprintf("https://api.%s.port.io", strings.ToLower(c.Region))
	}
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
