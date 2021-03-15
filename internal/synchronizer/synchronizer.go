package synchronizer

import (
	"context"
	"log"
	"time"

	"github.com/bynov/webhook-service/internal/batch"
	"github.com/bynov/webhook-service/internal/domain"
)

type WebhookRepository interface {
	GetLiteWebhooks(ctx context.Context, from, to time.Time) ([]domain.Webhook, error)
	SaveBatch(ctx context.Context, webhooks []domain.Webhook) error
}

type SlaveWebhookProvider interface {
	GetLiteWebhooks(ctx context.Context, from, to time.Time) ([]domain.Webhook, error)
	GetWebhooksByIDs(ctx context.Context, ids []string) ([]domain.Webhook, error)
}

type Synchronizer struct {
	webhookRepo          WebhookRepository
	slaveWebhookProvider SlaveWebhookProvider
	batchProvider        *batch.Batch

	errChan chan error
}

func New(webhookRepo WebhookRepository, slaveWebhookProvider SlaveWebhookProvider, batchProvider *batch.Batch) Synchronizer {
	return Synchronizer{
		webhookRepo:          webhookRepo,
		slaveWebhookProvider: slaveWebhookProvider,
		batchProvider:        batchProvider,

		errChan: make(chan error, 10),
	}
}

func (s Synchronizer) Run(syncInterval time.Duration) {
	for range time.NewTicker(syncInterval).C {
		now := time.Now().UTC()

		// Get webhooks from master from [-5; -1]
		masterWebhookSlice, err := s.webhookRepo.GetLiteWebhooks(
			context.Background(),
			now.Add(-5*time.Minute),
			now.Add(-1*time.Minute),
		)
		if err != nil {
			s.errChan <- err
			continue
		}

		// Convert master webhook slice to map
		var masterWebhooks = make(map[string][]time.Time)

		for _, v := range masterWebhookSlice {
			masterWebhooks[v.PayloadHash] = append(masterWebhooks[v.PayloadHash], v.RecievedAt)
		}

		log.Printf("master webhooks: %v", masterWebhooks)

		// Get webhooks from slave from [-4; -2]
		slaveWebhooks, err := s.slaveWebhookProvider.GetLiteWebhooks(
			context.Background(),
			now.Add(-4*time.Minute),
			now.Add(-2*time.Minute),
		)
		if err != nil {
			s.errChan <- err
			continue
		}

		// Find missing record IDs
		var missingIDs []string

		for i := range slaveWebhooks {
			var found bool

			// If not found - assume that it missing from master
			masterTimes, ok := masterWebhooks[slaveWebhooks[i].PayloadHash]
			if !ok {
				missingIDs = append(missingIDs, slaveWebhooks[i].ID)
				continue
			}

			// Webhook already presented - iterate over all times if not found - add it.
			for _, tm := range masterTimes {
				log.Println(slaveWebhooks[i].RecievedAt.Sub(tm))
				log.Println(slaveWebhooks[i].RecievedAt.Sub(tm) <= time.Minute)

				if slaveWebhooks[i].RecievedAt.Sub(tm) <= time.Minute {
					found = true
					break
				}
			}

			if !found {
				missingIDs = append(missingIDs, slaveWebhooks[i].ID)
			}
		}

		if len(missingIDs) == 0 {
			continue
		}

		// Get webhooks from slave by given IDs
		webhooksToInsert, err := s.slaveWebhookProvider.GetWebhooksByIDs(context.Background(), missingIDs)
		if err != nil {
			s.errChan <- err
			continue
		}

		// Add webhooks to batch provider
		err = s.batchProvider.Add(context.Background(), webhooksToInsert...)
		if err != nil {
			s.errChan <- err
			continue
		}
	}
}

func (s Synchronizer) Errors() <-chan error {
	return s.errChan
}
