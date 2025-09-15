package entity

import (
	"time"
)

// MessageReceipt represents message_receipts table
type MessageReceipt struct {
	ReceiptID string     `json:"receipt_id" db:"receipt_id"`
	AccountID string     `json:"account_id" db:"account_id"`
	MessageID string     `json:"message_id" db:"message_id"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateMessageReceiptRequest represents request to create message receipt
type CreateMessageReceiptRequest struct {
	ReceiptID string    `json:"receipt_id"`
	AccountID string    `json:"account_id"`
	MessageID string    `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}