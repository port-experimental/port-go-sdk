package webhooks

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/port-experimental/port-go-sdk/pkg/httpx"
)

// Post sends a JSON payload to a webhook URL with optional signature header.
func Post(ctx context.Context, hc httpx.Doer, url, secret string, payload any) error {
	if hc == nil {
		hc = httpx.New()
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	httpx.SetUserAgent(req, "")
	if secret != "" {
		req.Header.Set("X-Signature", sign(secret, body))
	}
	resp, err := httpx.DoWithRetry(ctx, hc, req, 3)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook POST failed: %s", resp.Status)
	}
	return nil
}

func sign(secret string, payload []byte) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write(payload)
	return hex.EncodeToString(m.Sum(nil))
}
