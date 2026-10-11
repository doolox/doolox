package cms

import (
	"runtime"

	"fyne.io/systray"
)

// TrayIcon holds the icon assets used for the system tray icon.
type TrayIcon struct {
	// PNG is used for the tray icon on macOS and Linux.
	PNG []byte
	// ICO is used for the tray icon on Windows.
	ICO []byte
}

// bytes returns the icon bytes that are appropriate for the current
// operating system.
func (t TrayIcon) bytes() []byte {
	if runtime.GOOS == "windows" {
		return t.ICO
	}

	return t.PNG
}

// runTray shows the Doolox tray icon and blocks until the tray is quit.
// onShutdown is called after the tray icon has been removed so the rest of
// the application can shut down cleanly.
func runTray(icon TrayIcon, onShutdown func()) {
	systray.Run(func() {
		systray.SetIcon(icon.bytes())
		systray.SetTitle("Doolox")
		systray.SetTooltip("Doolox is running")

		quit := systray.AddMenuItem("Quit", "Quit Doolox")

		go func() {
			<-quit.ClickedCh
			systray.Quit()
		}()
	}, func() {
		if onShutdown != nil {
			onShutdown()
		}
	})
}
