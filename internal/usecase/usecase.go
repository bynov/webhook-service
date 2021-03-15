package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bynov/webhook-service/internal/batch"
	"github.com/bynov/webhook-service/internal/domain"
)

// WebhookRepository is the approach to reach webhook internal repo.
type WebhookRepository interface {
	GetWebhooksByIDs(ctx context.Context, ids []string) ([]domain.Webhook, error)
	GetLiteWebhooks(ctx context.Context, from, to time.Time) ([]domain.Webhook, error)
}

type Usecase struct {
	batch       *batch.Batch
	webhookRepo WebhookRepository
	timeout     time.Duration
}

func NewUsecase(batch *batch.Batch, webhookRepo WebhookRepository, timeout time.Duration) Usecase {
	return Usecase{
		batch:       batch,
		webhookRepo: webhookRepo,
		timeout:     timeout,
	}
}

func (u *Usecase) AddWebhook(c context.Context, webhook domain.Webhook) error {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	if err := u.batch.Add(ctx, webhook); err != nil {
		return fmt.Errorf("failed to add webhook: %w", err)
	}

	return nil
}

func (u *Usecase) GetWebhooksByIDs(c context.Context, ids []string) ([]domain.Webhook, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	webhooks, err := u.webhookRepo.GetWebhooksByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks by ids: %w", err)
	}

	return webhooks, nil
}

func (u *Usecase) GetLiteWebhooks(c context.Context, from, to time.Time) ([]domain.Webhook, error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	webhooks, err := u.webhookRepo.GetLiteWebhooks(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get lite webhooks: %w", err)
	}

	return webhooks, nil
}
