package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/port-experimental/port-go-sdk/pkg/config"
)

func TestClientDo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Fatalf("missing auth header")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	cfg := config.Config{
		APIToken: "token",
	}
	cfg.BaseURL = srv.URL
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if err := c.Do(context.Background(), http.MethodGet, "/v1/test", nil, nil); err != nil {
		t.Fatalf("do failed: %v", err)
	}
}
