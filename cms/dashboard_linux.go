//go:build linux

package cms

import "log"

// launchDashboard opens the Doolox admin interface in the system browser. We
// deliberately do not use a native web view on Linux: embedding WebKitGTK
// inside a Go process is unreliable — JavaScriptCore's garbage collector
// signal conflicts with the Go runtime, and the web view cannot be re-created
// after its window is closed.
func launchDashboard() {
	if err := openBrowser(adminURL); err != nil {
		log.Printf("open dashboard: %v", err)
	}
}
