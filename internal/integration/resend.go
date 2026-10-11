package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Resend sends transactional email. If RESEND_API_KEY is empty, sending is disabled.
type ResendClient struct {
	APIKey string
	From   string
	Client *http.Client
}

func NewResendClient() *ResendClient {
	return &ResendClient{
		APIKey: strings.TrimSpace(os.Getenv("RESEND_API_KEY")),
		From:   strings.TrimSpace(os.Getenv("RESEND_FROM_EMAIL")),
		Client: &http.Client{Timeout: 10 * time.Second}}
}

func (r *ResendClient) Enabled() bool {
	return r != nil && r.APIKey != "" && r.From != ""
}

func (r *ResendClient) Send(ctx context.Context, to, subject, html string) error {
	if !r.Enabled() {
		return nil
	}

	body, err := json.Marshal(map[string]interface{}{"from": r.From, "to": []string{to}, "subject": subject, "html": html})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("resend returned HTTP %d", resp.StatusCode)
	}

	return nil
}
