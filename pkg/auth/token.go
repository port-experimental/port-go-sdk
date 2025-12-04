package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/httpx"
)

// TokenSource returns valid bearer tokens for Port API calls.
type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

// NewTokenSource chooses between static token or client credentials flow.
func NewTokenSource(cfg config.Config, hc httpx.Doer) TokenSource {
	if cfg.APIToken != "" {
		return &staticToken{token: cfg.APIToken}
	}
	if hc == nil {
		hc = httpx.New()
	}
	return &clientCredsSource{
		cfg: cfg,
		hc:  hc,
	}
}

type staticToken struct {
	token string
}

func (s *staticToken) Token(context.Context) (string, error) {
	return s.token, nil
}

type clientCredsSource struct {
	cfg config.Config
	hc  httpx.Doer

	mu      sync.Mutex
	token   string
	expires time.Time
}

func (c *clientCredsSource) Token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Until(c.expires) > 30*time.Second {
		return c.token, nil
	}
	if err := c.refresh(ctx); err != nil {
		return "", err
	}
	return c.token, nil
}

func (c *clientCredsSource) refresh(ctx context.Context) error {
	payload := map[string]string{
		"clientId":     c.cfg.ClientID,
		"clientSecret": c.cfg.ClientSecret,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseEndpoint()+"/v1/auth/access_token", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	httpx.SetUserAgent(req, "")
	resp, err := httpx.DoWithRetry(ctx, c.hc, req, 3)
	if err != nil {
		return fmt.Errorf("port auth: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("port auth failed: %s", resp.Status)
	}
	var out struct {
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expiresIn"`
	}
	var raw bytes.Buffer
	if err := json.NewDecoder(io.TeeReader(resp.Body, &raw)).Decode(&out); err != nil {
		return err
	}
	if out.AccessToken == "" {
		return fmt.Errorf("port auth: empty token response=%s", raw.String())
	}
	c.token = out.AccessToken
	if out.ExpiresIn == 0 {
		out.ExpiresIn = 3600
	}
	c.expires = time.Now().Add(time.Duration(out.ExpiresIn) * time.Second)
	return nil
}
