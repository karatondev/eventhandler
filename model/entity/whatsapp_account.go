package entity

import (
	"time"

	sharedmodel "zaplio/shared/model"
)

// WhatsAppAccount represents whatsapp_accounts table
type WhatsAppAccount struct {
	AccountID      string               `json:"account_id" db:"account_id"`
	UserID         string               `json:"user_id" db:"user_id"`
	AccountName    string               `json:"account_name" db:"account_name"`
	AccountAlias   *string              `json:"account_alias" db:"account_alias"`
	PhoneNumber    *string              `json:"phone_number" db:"phone_number"`
	SenderJID      *string              `json:"sender_jid" db:"sender_jid"`
	SessionData    []byte               `json:"session_data" db:"session_data"`
	ConnectStatus  sharedmodel.EventType `json:"connect_status" db:"connect_status"`
	IsActive       bool                 `json:"is_active" db:"is_active"`
	InitiatedAt    *time.Time           `json:"initiated_at" db:"initiated_at"`
	ConnectedAt    *time.Time           `json:"connected_at" db:"connected_at"`
	DisconnectedAt *time.Time           `json:"disconnected_at" db:"disconnected_at"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
	CreatedBy      string               `json:"created_by" db:"created_by"`
	UpdatedAt      *time.Time           `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time           `json:"deleted_at" db:"deleted_at"`
}

type WhatsAppAccountReq struct {
	AccountID      string               `json:"account_id"`
	AccountName    string               `json:"account_name"`
	PhoneNumber    *string              `json:"phone_number"`
	SenderJID      *string              `json:"sender_jid"`
	SessionData    []byte               `json:"session_data"`
	ConnectStatus  sharedmodel.EventType `json:"connect_status"`
	InitiatedAt    *time.Time           `json:"initiated_at" db:"initiated_at"`
	ConnectedAt    *time.Time           `json:"connected_at" db:"connected_at"`
	DisconnectedAt *time.Time           `json:"disconnected_at" db:"disconnected_at"`
}
