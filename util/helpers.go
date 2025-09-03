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
// Example: "a81bc81b-dead-4e5d-abff-90865d1e13b1@s.whatsapp.net" -> "a81bc81b-dead-4e5d-abff-90865d1e13b1"
func ExtractJIDPrefix(jid string) string {
	if idx := strings.Index(jid, "@"); idx != -1 {
		return jid[:idx]
	}
	return jid // Return original if no @ found
}
