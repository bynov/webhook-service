package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bynov/webhook-service/internal/domain"
)

type WebhookRepo interface {
	SaveBatch(ctx context.Context, webhooks []domain.Webhook) error
}

type Batch struct {
	maxDataBeforeFlush int

	data      []domain.Webhook
	dataMutex sync.Mutex

	webhookRepo WebhookRepo
	ticker      *time.Ticker

	errChan chan error
}

func New(webhookRepo WebhookRepo, flushInterval time.Duration, maxDataBeforeFlush int) *Batch {
	return &Batch{
		webhookRepo:        webhookRepo,
		ticker:             time.NewTicker(flushInterval),
		maxDataBeforeFlush: maxDataBeforeFlush,

		errChan: make(chan error, 10),
	}
}

func (b *Batch) Add(webhook domain.Webhook) error {
	b.dataMutex.Lock()
	defer b.dataMutex.Unlock()

	b.data = append(b.data, webhook)

	if len(b.data) >= b.maxDataBeforeFlush {
		if err := b.flush(); err != nil {
			return fmt.Errorf("batch: flush error: %w", err)
		}
	}

	return nil
}

func (b *Batch) Start() {
	for range b.ticker.C {
		if err := b.Flush(); err != nil {
			b.errChan <- fmt.Errorf("batch: flush error: %w", err)
		}
	}
}

func (b *Batch) Flush() error {
	b.dataMutex.Lock()

	if len(b.data) == 0 {
		return nil
	}

	var data = make([]domain.Webhook, len(b.data))
	copy(data, b.data)

	b.dataMutex.Unlock()

	// TODO: err & context
	return b.webhookRepo.SaveBatch(context.TODO(), data)
}

func (b *Batch) Errors() <-chan error {
	return b.errChan
}

func (b *Batch) flush() error {
	// TODO: err & context
	return b.webhookRepo.SaveBatch(context.TODO(), b.data)
}
