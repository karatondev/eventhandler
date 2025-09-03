package entity

import (
	"time"
)

// AccountEvent represents account_events table
type AccountEvent struct {
	EventID   string     `json:"event_id" db:"event_id"`
	AccountID string     `json:"account_id" db:"account_id"`
	EventType string     `json:"event_type" db:"event_type"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Data      []byte     `json:"data" db:"data"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateAccountEventRequest represents request to create account event
type CreateAccountEventRequest struct {
	EventID   string    `json:"event_id"`
	AccountID string    `json:"account_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}
