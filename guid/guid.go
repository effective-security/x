package guid

import (
	"uuid"
)

// MustCreate returns a UUID string (currently UUID v4, lowercase).
//
// Deprecated: use uuid.NewV7 or uuid.NewV4 from the standard library uuid package.
func MustCreate() string {
	return uuid.NewV4().String()
}
