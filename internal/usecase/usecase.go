package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bynov/webhook-service/internal/batch"
	"github.com/bynov/webhook-service/internal/domain"
)

type Usecase struct {
	batch   *batch.Batch
	timeout time.Duration
}

func NewUsecase(batch *batch.Batch, timeout time.Duration) Usecase {
	return Usecase{
		batch:   batch,
		timeout: timeout,
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
