package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/port-experimental/port-go-sdk/pkg/auth"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/httpx"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
)

// Client is the root Port API client.
type Client struct {
	baseURL     string
	hc          httpx.Doer
	tokenSource auth.TokenSource
	userAgent   string
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

// New constructs a Port API client.
func New(cfg config.Config, opts ...Option) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	c := &Client{
		baseURL:   strings.TrimRight(cfg.BaseEndpoint(), "/"),
		hc:        httpx.New(),
		userAgent: "port-go-sdk/0.1",
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.hc == nil {
		c.hc = httpx.New()
	}
	c.tokenSource = auth.NewTokenSource(cfg, c.hc)
	return c, nil
}

// do issues an HTTP request and decodes JSON result into out (if non-nil).
// Do is exported so service packages can invoke Port API endpoints.
func (c *Client) Do(ctx context.Context, method, path string, body any, out any) error {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	}
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
	resp, err := httpx.DoWithRetry(ctx, c.hc, req, 3)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		payload, _ := io.ReadAll(resp.Body)
		return &porter.Error{StatusCode: resp.StatusCode, Message: string(payload), Body: payload}
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// ping ensures credentials are valid.
func (c *Client) ping(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/health", nil)
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
