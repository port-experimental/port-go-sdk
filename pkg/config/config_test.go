package config

import (
	"os"
	"testing"
)

func TestBaseEndpoint(t *testing.T) {
	c := Config{}
	if got := c.BaseEndpoint(); got != "https://api.port.io" {
		t.Fatalf("expected default EU base, got %s", got)
	}
	c.Region = "US"
	if got := c.BaseEndpoint(); got != "https://api.us.port.io" {
		t.Fatalf("expected US base, got %s", got)
	}
	c.BaseURL = "https://custom.example.com/v1"
	if got := c.BaseEndpoint(); got != "https://custom.example.com/v1" {
		t.Fatalf("expected override to win, got %s", got)
	}
}

func TestLoadPrefersToken(t *testing.T) {
	t.Setenv("PORT_ACCESS_TOKEN", "abc")
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.APIToken != "abc" {
		t.Fatalf("expected direct token")
	}
}

func TestLoadRequiresCreds(t *testing.T) {
	os.Unsetenv("PORT_ACCESS_TOKEN")
	os.Unsetenv("PORT_CLIENT_ID")
	os.Unsetenv("PORT_CLIENT_SECRET")
	if _, err := Load(""); err == nil {
		t.Fatalf("expected validation error without creds")
	}
	t.Setenv("PORT_CLIENT_ID", "id")
	t.Setenv("PORT_CLIENT_SECRET", "secret")
	if _, err := Load(""); err != nil {
		t.Fatalf("unexpected error once creds set: %v", err)
	}
}
