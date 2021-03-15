package domain

import "time"

type Webhook struct {
	ID          string    `json:"id"`
	Payload     string    `json:"payload"`
	PayloadHash string    `json:"payload_hash"`
	RecievedAt  time.Time `json:"recieved_at"`
}
