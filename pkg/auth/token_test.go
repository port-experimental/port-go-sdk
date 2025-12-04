package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/httpx"
)

func TestStaticToken(t *testing.T) {
	cfg := config.Config{APIToken: "abc"}
	src := NewTokenSource(cfg, nil)
	tok, err := src.Token(context.Background())
	if err != nil || tok != "abc" {
		t.Fatalf("expected token, got %s err %v", tok, err)
	}
}

func TestClientCredentialsSource(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"accessToken": "token123",
			"expiresIn":   1,
		})
	}))
	defer srv.Close()
	cfg := config.Config{
		ClientID:     "id",
		ClientSecret: "secret",
	}
	cfg.BaseURL = srv.URL
	src := NewTokenSource(cfg, httpx.New())
	ctx := context.Background()
	tok, err := src.Token(ctx)
	if err != nil || tok != "token123" {
		t.Fatalf("unexpected token: %s err %v", tok, err)
	}
	// wait for expiry to force refresh
	time.Sleep(2 * time.Second)
	tok, err = src.Token(ctx)
	if err != nil || tok != "token123" {
		t.Fatalf("expected refresh success: %s err %v", tok, err)
	}
}
