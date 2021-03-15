package domain

import "time"

// Webhook represents webhook struct.
type Webhook struct {
	ID          string    `json:"id"`
	Payload     string    `json:"payload"`
	PayloadHash string    `json:"payload_hash"`
	RecievedAt  time.Time `json:"recieved_at"`
}
