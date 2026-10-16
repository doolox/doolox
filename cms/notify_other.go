//go:build !linux && !windows

package cms

import "log"

// notify is a no-op on platforms without a native notification backend; the
// message is just logged.
func notify(title, message string) {
	log.Printf("[notify] %s: %s", title, message)
}
