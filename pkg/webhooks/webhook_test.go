package webhooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	var gotSig string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	err := Post(context.Background(), srv.Client(), srv.URL, "secret", map[string]string{"hello": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotSig == "" {
		t.Fatalf("expected signature header")
	}
}
