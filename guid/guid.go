package guid

import (
	"crypto/rand"
	"fmt"
)

// randRead is the function used to read random bytes; it can be overridden for testing.
var randRead = rand.Read

// MustCreate returns GUID
func MustCreate() string {
	b := make([]byte, 16)
	_, err := randRead(b)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
