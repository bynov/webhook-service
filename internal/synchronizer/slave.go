package synchronizer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bynov/webhook-service/internal/domain"
)

const (
	getLiteWebhooks  = "/webhooks/lite"
	getWebhooksByIDs = "/webhooks"
)

// TODO: can't find better name :(
type WebhookProviderFromSlave struct {
	client *http.Client
	url    string
}

func NewWebhookProviderFromSlave(url string, timeout time.Duration) WebhookProviderFromSlave {
	return WebhookProviderFromSlave{
		client: &http.Client{
			Timeout: timeout,
		},
		url: url,
	}
}

// GetLiteWebhooks makes a http call to slave to get lite webhooks.
func (w WebhookProviderFromSlave) GetLiteWebhooks(ctx context.Context, from, to time.Time) ([]domain.Webhook, error) {
	values := url.Values{
		"from": {strconv.FormatInt(from.UTC().Unix(), 10)},
		"to":   {strconv.FormatInt(to.UTC().Unix(), 10)},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.url+getLiteWebhooks+"?"+values, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request in get lite webhooks from slave: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get lite webhooks from slave: %w", err)
	}

	defer resp.Body.Close()

	var webhooks []domain.Webhook

	if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return nil, fmt.Errorf("failed to decode in get lite webhooks from slave: %w", err)
	}

	return webhooks, nil
}

// GetWebhooksByIDs makes a http call to slave to get webhooks by ids.
func (w WebhookProviderFromSlave) GetWebhooksByIDs(ctx context.Context, ids []string) ([]domain.Webhook, error) {
	values := url.Values{
		"ids[]": ids,
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.url+getWebhooksByIDs+"?"+values, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request in get webhooks by ids from slave: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks by ids from slave: %w", err)
	}

	defer resp.Body.Close()

	var webhooks []domain.Webhook

	if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return nil, fmt.Errorf("failed to decode in get webhooks by ids from slave: %w", err)
	}

	return webhooks, nil
}
