package entity

import (
	"time"
)

// MessageInbound represents message_inbounds table
type MessageInbound struct {
	InboundID   string     `json:"inbound_id" db:"inbound_id"`
	AccountID   string     `json:"account_id" db:"account_id"`
	FromMe      *bool      `json:"from_me" db:"from_me"`
	MessageID   string     `json:"message_id" db:"message_id"`
	Sender      string     `json:"sender" db:"sender"`
	MessageType string     `json:"message_type" db:"message_type"`
	ReceivedAt  time.Time  `json:"received_at" db:"received_at"`
	Data        []byte     `json:"data" db:"data"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateMessageInboundRequest represents request to create message inbound
type CreateMessageInboundRequest struct {
	EventID     string    `json:"event_id"`
	AccountID   string    `json:"account_id"`
	FromMe      *bool     `json:"from_me"`
	MessageID   string    `json:"message_id"`
	Sender      string    `json:"sender"`
	MessageType string    `json:"message_type"`
	ReceivedAt  time.Time `json:"received_at"`
	Data        []byte    `json:"data"`
}
