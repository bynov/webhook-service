package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/bynov/webhook-service/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type webhookBatch struct {
	payloads      []string
	payloadHashes []string
	recievedAt    []time.Time
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{pool: pool}
}

func (p PostgresRepository) SaveBatch(ctx context.Context, webhooks []domain.Webhook) error {
	batch := toWebhookBatch(webhooks)

	_, err := p.pool.Exec(
		ctx,
		`INSERT INTO "webhooks" (
				"payload",
				"payload_hash",
				"received_at"
			)
			SELECT
				unnest($1::varchar[]),
				unnest($2::char[]),
				unnest($3::timestamp[])`,
		batch.payloads,
		batch.payloadHashes,
		batch.recievedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert batch, got: %w", err)
	}

	return nil
}

func toWebhookBatch(webhooks []domain.Webhook) webhookBatch {
	var out = webhookBatch{
		payloads:      make([]string, len(webhooks)),
		payloadHashes: make([]string, len(webhooks)),
		recievedAt:    make([]time.Time, len(webhooks)),
	}

	for i := range webhooks {
		out.payloads[i] = webhooks[i].Payload
		out.payloadHashes[i] = webhooks[i].PayloadHash
		out.recievedAt[i] = webhooks[i].RecievedAt
	}

	return out
}
