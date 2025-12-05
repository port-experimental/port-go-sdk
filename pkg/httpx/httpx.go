// Package httpx provides HTTP client utilities including retry logic,
// request cloning, and connection pooling for the Port API SDK.
package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultUserAgent = "port-go-sdk/0.1"

// rng is a package-level random number generator for jitter calculations.
// Using a local rand.Rand instance instead of the deprecated global rand.Seed.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Doer matches http.Client.Do.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// New returns an http.Client with sensible defaults.
func New() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 32,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
}

// SetUserAgent sets a default UA if not provided.
func SetUserAgent(req *http.Request, ua string) {
	if req == nil {
		return
	}
	if strings.TrimSpace(ua) == "" {
		ua = defaultUserAgent
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", ua)
	}
}

// DoWithRetry retries on 5xx/429 with exponential backoff + jitter.
func DoWithRetry(ctx context.Context, client Doer, req *http.Request, attempts int) (*http.Response, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var err error
	var resp *http.Response
	for i := 1; i <= attempts; i++ {
		cloned, cerr := cloneRequest(req)
		if cerr != nil {
			return nil, cerr
		}
		resp, err = client.Do(cloned)
		if err == nil && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}
		var wait time.Duration
		if resp != nil {
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if seconds, parseErr := strconv.Atoi(ra); parseErr == nil {
					wait = time.Duration(seconds) * time.Second
				}
			}
			_ = resp.Body.Close()
		}
		if wait == 0 {
			backoff := time.Duration(math.Pow(2, float64(i-1))) * time.Second
			jitter := time.Duration(rng.Intn(500)) * time.Millisecond
			wait = backoff + jitter
		}
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if err == nil && resp != nil {
		return resp, fmt.Errorf("max retries exceeded: %s", resp.Status)
	}
	return nil, err
}

// cloneRequest copies the request and body so it can be retried.
func cloneRequest(req *http.Request) (*http.Request, error) {
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	cloned := req.Clone(req.Context())
	if req.Body == nil || req.Body == http.NoBody {
		return cloned, nil
	}
	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()
	defer buffers.Put(buf)
	if _, err := io.Copy(buf, req.Body); err != nil {
		return nil, err
	}
	_ = req.Body.Close()
	bodyCopy := bytes.NewReader(buf.Bytes())
	req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	cloned.Body = io.NopCloser(bodyCopy)
	if req.ContentLength > 0 {
		cloned.ContentLength = int64(buf.Len())
		req.ContentLength = int64(buf.Len())
	}
	return cloned, nil
}

var buffers = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}
