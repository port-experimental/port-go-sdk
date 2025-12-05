// Package client provides the main Port API client and service accessors.
// It handles authentication, HTTP requests, retries, and provides access to
// all Port API services (entities, blueprints, automations, etc.).
//
// Example usage:
//
//	cfg, _ := config.Load(".env")
//	cli, _ := client.New(cfg)
//	defer cli.Close()
//	entity, _ := cli.Entities().Get(ctx, "blueprint", "identifier")
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/auth"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/httpx"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
	"github.com/port-experimental/port-go-sdk/pkg/version"
)

// Client is the root Port API client.
type Client struct {
	baseURL       string
	hc            httpx.Doer
	tokenSource   auth.TokenSource
	userAgent     string
	verbose       bool
	logger        *log.Logger
	logFile       *os.File // Track file handle for cleanup
	respLimit     int64
	bufferPool    sync.Pool
	retryAttempts int // Number of retry attempts for failed requests
}

// Option mutates the Client.
type Option func(*Client)

// WithHTTPClient overrides the HTTP client.
func WithHTTPClient(h httpx.Doer) Option {
	return func(c *Client) { c.hc = h }
}

// WithUserAgent sets a custom user agent.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithRetryAttempts sets the number of retry attempts for failed requests.
// Default is 3. Set to 1 to disable retries.
func WithRetryAttempts(attempts int) Option {
	return func(c *Client) {
		if attempts < 1 {
			attempts = 1
		}
		c.retryAttempts = attempts
	}
}

// New constructs a Port API client from the provided configuration.
// The client handles authentication automatically using either an API token
// or client credentials from the config.
//
// Example:
//
//	cfg, _ := config.Load(".env")
//	cli, err := client.New(cfg, client.WithUserAgent("my-app/1.0"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
func New(cfg config.Config, opts ...Option) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	c := &Client{
		baseURL:       strings.TrimRight(cfg.BaseEndpoint(), "/"),
		hc:            httpx.New(),
		userAgent:     version.UserAgent(),
		respLimit:     10 << 20,
		retryAttempts: 3, // Default to 3 retry attempts
		bufferPool: sync.Pool{
			New: func() any { return new(bytes.Buffer) },
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.hc == nil {
		c.hc = httpx.New()
	}
	c.tokenSource = auth.NewTokenSource(cfg, c.hc)
	c.initVerboseLogger()
	c.initResponseLimit()
	return c, nil
}

// Do issues an HTTP request to the Port API and decodes the JSON response
// into out (if non-nil). It handles authentication, retries, and error handling.
// This method is exported so service packages can invoke Port API endpoints.
//
// The context controls the request lifetime. If the context is canceled or
// times out, the request will be aborted.
//
// If out is nil, the response body is discarded. Otherwise, it must be a pointer
// to a struct that can be unmarshaled from JSON.
func (c *Client) Do(ctx context.Context, method, path string, body any, out any) error {
	var rdr io.Reader
	cleanup := func() {}
	if body != nil {
		encReader, release, err := c.encodeBody(body)
		if err != nil {
			return err
		}
		rdr = encReader
		cleanup = release
	}
	defer cleanup()

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, rdr)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	httpx.SetUserAgent(req, c.userAgent)
	token, err := c.tokenSource.Token(ctx)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	start := time.Now()
	c.verbosef("--> %s %s", method, path)
	resp, err := httpx.DoWithRetry(ctx, c.hc, req, c.retryAttempts)
	if err != nil {
		c.verbosef("<!! %s %s error=%v", method, path, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		payload, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			// If we can't read the error body, still return the HTTP error
			payload = []byte(fmt.Sprintf("failed to read error body: %v", readErr))
		}
		err := &porter.Error{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("port api: %s %s", resp.Status, path),
			Body:       payload,
		}
		c.verbosef("<-- %s %s status=%d duration=%s", method, path, resp.StatusCode, time.Since(start))
		return err
	}
	if out != nil {
		reader := io.Reader(resp.Body)
		if c.respLimit > 0 {
			reader = io.LimitReader(resp.Body, c.respLimit)
		}
		err := json.NewDecoder(reader).Decode(out)
		c.verbosef("<-- %s %s status=%d duration=%s", method, path, resp.StatusCode, time.Since(start))
		return err
	}
	c.verbosef("<-- %s %s status=%d duration=%s", method, path, resp.StatusCode, time.Since(start))
	return nil
}

// Ping ensures credentials are valid by checking the health endpoint.
// This method is exported for users who want to verify their credentials.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/health", http.NoBody)
	if err != nil {
		return fmt.Errorf("ping: failed to create request: %w", err)
	}
	httpx.SetUserAgent(req, c.userAgent)
	token, err := c.tokenSource.Token(ctx)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpx.DoWithRetry(ctx, c.hc, req, 1)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("ping failed: %s", resp.Status)
	}
	return nil
}

func (c *Client) initVerboseLogger() {
	raw, ok := os.LookupEnv("PORT_SDK_VERBOSE")
	if !ok {
		return
	}
	enabled, err := strconv.ParseBool(raw)
	if err != nil || !enabled {
		return
	}
	logPath := os.Getenv("PORT_SDK_VERBOSE_FILE")
	var w io.Writer = os.Stdout
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err == nil {
			c.logFile = f
			w = f
		}
	}
	c.logger = log.New(w, "[port-sdk] ", log.LstdFlags|log.Lmicroseconds)
	c.verbose = true
}

// Close releases resources held by the client, including log file handles.
func (c *Client) Close() error {
	if c.logFile != nil {
		return c.logFile.Close()
	}
	return nil
}

func (c *Client) verbosef(format string, args ...any) {
	if !c.verbose || c.logger == nil {
		return
	}
	// Sanitize arguments to prevent logging sensitive data
	sanitized := make([]any, len(args))
	for i, arg := range args {
		sanitized[i] = c.sanitizeLogArg(arg)
	}
	c.logger.Printf(format, sanitized...)
}

// sanitizeLogArg removes sensitive information from log arguments.
func (c *Client) sanitizeLogArg(arg any) any {
	if arg == nil {
		return nil
	}
	switch v := arg.(type) {
	case string:
		// Check for common sensitive patterns
		if strings.HasPrefix(v, "Bearer ") && len(v) > 20 {
			return "Bearer [REDACTED]"
		}
		if strings.Contains(strings.ToLower(v), "password") ||
			strings.Contains(strings.ToLower(v), "secret") ||
			strings.Contains(strings.ToLower(v), "token") ||
			strings.Contains(strings.ToLower(v), "key") {
			// Don't log potentially sensitive strings
			return "[REDACTED]"
		}
		return v
	case error:
		// Error messages might contain sensitive data
		errStr := v.Error()
		if strings.Contains(errStr, "Bearer ") ||
			strings.Contains(strings.ToLower(errStr), "secret") ||
			strings.Contains(strings.ToLower(errStr), "token") {
			return fmt.Errorf("[REDACTED: error may contain sensitive data]")
		}
		return v
	default:
		return arg
	}
}

func (c *Client) initResponseLimit() {
	if raw := os.Getenv("PORT_SDK_MAX_RESPONSE_BYTES"); raw != "" {
		if v, err := strconv.ParseInt(raw, 10, 64); err == nil && v > 0 {
			c.respLimit = v
		}
	}
}

func (c *Client) encodeBody(body any) (io.Reader, func(), error) {
	if rdr, ok := body.(io.Reader); ok {
		return rdr, func() {}, nil
	}
	buf := c.getBuffer()
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		c.putBuffer(buf)
		return nil, func() {}, err
	}
	reader := bytes.NewReader(buf.Bytes())
	return reader, func() {
		buf.Reset()
		c.putBuffer(buf)
	}, nil
}

func (c *Client) getBuffer() *bytes.Buffer {
	raw := c.bufferPool.Get()
	buf, ok := raw.(*bytes.Buffer)
	if !ok || buf == nil {
		buf = &bytes.Buffer{}
	}
	buf.Reset()
	return buf
}

func (c *Client) putBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	c.bufferPool.Put(buf)
}
