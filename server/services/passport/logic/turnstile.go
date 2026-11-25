package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"aqua/pkg/infra/models"
	"aqua/pkg/transport/client"
)

// verifyTurnstileToken validates the provided token against Cloudflare's endpoint.
func verifyTurnstileToken(ctx context.Context, cfg models.TurnstileConfig, token string, ip string) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if token == "" {
		return false, nil
	}
	if cfg.SecretKey == "" || cfg.VerifyURL == "" {
		return false, fmt.Errorf("turnstile config is not configured")
	}

	data := url.Values{}
	data.Set("secret", cfg.SecretKey)
	data.Set("response", token)
	if ip != "" {
		data.Set("remoteip", ip)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.VerifyURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := client.NewHTTPClient(10 * time.Second)
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result struct {
		Success     bool     `json:"success"`
		ChallengeTS string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes"`
		Action      string   `json:"action"`
		CData       string   `json:"cdata"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result.Success, nil
}
