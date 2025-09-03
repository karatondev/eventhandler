package util

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// MarshalToJSON marshals data to JSON []byte, returns error if failed
func MarshalToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// NowUTC returns current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}

// CreateDirectory create multiple directory.
func CreateDirectory(paths ...string) (err error) {
	for _, path := range paths {
		_, notExistError := os.Stat(path)
		if os.IsNotExist(notExistError) {
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return
}

// ExtractJIDPrefix extracts the part before @ from a JID string
// If the JID contains phone number format (digits:digits), returns the original JID
// If the JID contains UUID format, returns only the UUID part
// Examples:
// - "6289513689096:50@s.whatsapp.net" -> "6289513689096:50@s.whatsapp.net" (no extraction)
// - "a81bc81b-dead-4e5d-abff-90865d1e13b1@s.whatsapp.net" -> "a81bc81b-dead-4e5d-abff-90865d1e13b1"
func ExtractJIDPrefix(jid string) string {
	// Check if JID contains phone number format (digits:digits@domain)
	if strings.Contains(jid, ":") && strings.Contains(jid, "@") {
		// It's phone number format, return original JID
		return jid
	}

	// Extract UUID part (everything before @)
	if idx := strings.Index(jid, "@"); idx != -1 {
		return jid[:idx]
	}
	return jid // Return original if no @ found
}
