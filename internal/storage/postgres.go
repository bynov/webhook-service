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

// SaveBatch is used to save webhooks in batch request.
// It could seems a bit strange but it more faster than pg.Batch{}.
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
				unnest($2::char(40)[]),
				unnest($3::timestamp[]);`,
		batch.payloads,
		batch.payloadHashes,
		batch.recievedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert batch, got: %w", err)
	}

	return nil
}

// GetLiteWebhooks returns list of webhooks without full payload. Hahs provided instead.
func (p PostgresRepository) GetLiteWebhooks(ctx context.Context, from, to time.Time) ([]domain.Webhook, error) {
	rows, err := p.pool.Query(
		ctx,
		`SELECT
				id,
				payload_hash,
				received_at
			FROM
				"webhooks"
			WHERE $1 <= received_at AND received_at < $2;`,
		from,
		to,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get lite webhooks: %w", err)
	}

	defer rows.Close()

	var out []domain.Webhook

	for rows.Next() {
		var item domain.Webhook

		if err := rows.Scan(
			&item.ID,
			&item.PayloadHash,
			&item.RecievedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan in get lite webhooks: %w", err)
		}

		out = append(out, item)
	}

	return out, nil
}

// GetWebhooksByIDs returns list of full webhooks by ids.
func (p PostgresRepository) GetWebhooksByIDs(ctx context.Context, ids []string) ([]domain.Webhook, error) {
	rows, err := p.pool.Query(
		ctx,
		`SELECT
				id,
				payload,
				payload_hash,
				received_at
			FROM
				"webhooks"
			WHERE id = ANY($1);`,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks by ids: %w", err)
	}

	defer rows.Close()

	var out []domain.Webhook

	for rows.Next() {
		var item domain.Webhook

		if err := rows.Scan(
			&item.ID,
			&item.Payload,
			&item.PayloadHash,
			&item.RecievedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan in get webhooks by ids: %w", err)
		}

		out = append(out, item)
	}

	return out, nil
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
