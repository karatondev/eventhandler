package model

import (
	"encoding/json"
)

// MessageData represents a unified message data structure
// that combines all possible message properties without sender and metadata
type MessageData struct {
	// Common message properties
	Content string `json:"content,omitempty"`
	Caption string `json:"caption,omitempty"`

	// Media properties (Image, Audio, Video, Document)
	MimeType string `json:"mime_type,omitempty"`
	FileSize uint64 `json:"file_size,omitempty"`
	FileURL  string `json:"file_url,omitempty"`

	// Audio/Video specific
	Duration uint32 `json:"duration,omitempty"` // in seconds

	// Document specific
	FileName string `json:"file_name,omitempty"`

	// Location specific
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`

	// Reaction specific
	Text         string `json:"text,omitempty"`
	TargetKey    string `json:"target_key,omitempty"`
	TargetSender string `json:"target_sender,omitempty"`

	// Button response specific
	SelectedButtonID string `json:"selected_button_id,omitempty"`
	DisplayText      string `json:"display_text,omitempty"`

	// List response specific
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	SelectedRowID string `json:"selected_row_id,omitempty"`
}

// ToBytes converts MessageData to JSON bytes
func (md *MessageData) ToBytes() ([]byte, error) {
	return json.Marshal(md)
}

// MustToBytes converts MessageData to JSON bytes, panics on error
func (md *MessageData) MustToBytes() []byte {
	data, err := json.Marshal(md)
	if err != nil {
		panic(err)
	}
	return data
}
