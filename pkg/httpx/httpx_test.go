package httpx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoWithRetrySuccess(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&hits, 1) < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "hello" {
			t.Fatalf("expected body replay, got %q", body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	req, _ := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader("hello"))
	resp, err := DoWithRetry(context.Background(), srv.Client(), req, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected success status")
	}
}

type errDoer struct{}

func (errDoer) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("boom")
}

func TestDoWithRetryStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if _, err := DoWithRetry(ctx, errDoer{}, req, 5); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context cancel, got %v", err)
	}
}
