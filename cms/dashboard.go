package cms

import (
	"os"

	webview "github.com/webview/webview_go"
)

// launchDashboard opens the Doolox admin interface in a dedicated native web
// view window instead of the default browser. It blocks until the window is
// closed, so callers should run it in its own goroutine.
func launchDashboard() {
	// JavaScriptCore (WebKitGTK's JS engine) installs a signal handler for its
	// garbage collector. The Go runtime already installs handlers for nearly
	// every signal (including JSC's default SIGUSR1), which makes JSC abort
	// with "Overriding existing handler for signal 10". SIGCONT (18) is one of
	// the few signals the Go runtime leaves at its default disposition, so
	// point JSC at it. Respect an explicit user override.
	if os.Getenv("JSC_SIGNAL_FOR_GC") == "" {
		os.Setenv("JSC_SIGNAL_FOR_GC", "18")
	}

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Doolox Dashboard")
	w.SetSize(1100, 720, webview.HintNone)
	w.Navigate(adminURL)
	w.Run()
}
