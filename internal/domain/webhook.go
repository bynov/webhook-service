package domain

import "time"

type Webhook struct {
	Payload     string
	PayloadHash string
	RecievedAt  time.Time
}
