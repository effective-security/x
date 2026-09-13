package guid

import (
	"uuid"
)

// MustCreate returns GUID
// Deprecated: Usestandard UUID instead: uid.NewV7() or uuid.NewV4()
func MustCreate() string {
	return uuid.NewV7().String()
}
