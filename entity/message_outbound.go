package entity

import (
	"time"
)

// MessageOutbound represents message_outbounds table
type MessageOutbound struct {
	OutboundID  string     `json:"outbound_id" db:"outbound_id"`
	AccountID   string     `json:"account_id" db:"account_id"`
	MessageID   string     `json:"message_id" db:"message_id"`
	Recipient   string     `json:"recipient" db:"recipient"`
	MessageType string     `json:"message_type" db:"message_type"`
	SentAt      time.Time  `json:"sent_at" db:"sent_at"`
	Data        []byte     `json:"data" db:"data"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateMessageOutboundRequest represents request to create message outbound
type CreateMessageOutboundRequest struct {
	EventID     string    `json:"event_id"`
	AccountID   string    `json:"account_id"`
	MessageID   string    `json:"message_id"`
	Recipient   string    `json:"recipient"`
	MessageType string    `json:"message_type"`
	SentAt      time.Time `json:"sent_at"`
	Data        []byte    `json:"data"`
}
