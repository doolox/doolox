package cms

import webview "github.com/webview/webview_go"

// launchDashboard opens the Doolox admin interface in a dedicated native web
// view window instead of the default browser. It blocks until the window is
// closed, so callers should run it in its own goroutine.
func launchDashboard() {
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Doolox Dashboard")
	w.SetSize(1100, 720, webview.HintNone)
	w.Navigate(adminURL)
	w.Run()
}
