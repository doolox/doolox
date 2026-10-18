package cms

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os/exec"
	"runtime"

	"fyne.io/systray"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// TrayIcon holds the icon assets used for the system tray icon.
type TrayIcon struct {
	// SVG is rasterized for the tray icon on Linux.
	SVG []byte
	// PNG is used for the tray icon on macOS and other non-Windows platforms.
	PNG []byte
	// ICO is used for the tray icon on Windows.
	ICO []byte
}

// bytes returns the icon bytes that are appropriate for the current
// operating system.
func (t TrayIcon) bytes() ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return t.ICO, nil
	case "linux":
		data, err := renderSVGIcon(t.SVG)
		if err != nil {
			return nil, fmt.Errorf("render Linux tray icon: %w", err)
		}
		return data, nil
	default:
		return t.PNG, nil
	}
}

// adminURL and websiteURL are the local addresses of the admin interface and
// public website, respectively.
const (
	adminURL   = "http://localhost:5555"
	websiteURL = "http://localhost:5000"
)

// runTray shows the Doolox tray icon and blocks until the tray is quit.
// onShutdown is called after the tray icon has been removed so the rest of
// the application can shut down cleanly.
func runTray(icon TrayIcon, onShutdown func()) {
	iconData, err := icon.bytes()
	if err != nil {
		log.Fatalf("prepare tray icon: %v", err)
	}

	systray.Run(func() {
		systray.SetIcon(iconData)
		systray.SetTitle("Doolox")
		systray.SetTooltip("Doolox is running")

		systray.SetOnTapped(func() {
			if err := openBrowser(websiteURL); err != nil {
				log.Printf("open browser: %v", err)
			}
		})

		mGenerate := systray.AddMenuItem("Generate (ctrl+g)", "Generate the static site")
		mPublish := systray.AddMenuItem("Publish (ctrl+p)", "Generate the static site and push it with git")
		mDashboard := systray.AddMenuItem("Doolox Dashboard (ctrl+d)", "Open the Doolox admin interface")

		go func() {
			for {
				select {
				case <-mGenerate.ClickedCh:
					generateSite()
				case <-mPublish.ClickedCh:
					publishSite()
				case <-mDashboard.ClickedCh:
					go launchDashboard()
				}
			}
		}()

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

func renderSVGIcon(svg []byte) ([]byte, error) {
	const size = 64

	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg), oksvg.StrictErrorMode)
	if err != nil {
		return nil, err
	}

	canvas := image.NewRGBA(image.Rect(0, 0, size, size))
	icon.SetTarget(0, 0, size, size)
	dasher := rasterx.NewDasher(size, size, rasterx.NewScannerGV(size, size, canvas, canvas.Bounds()))
	icon.Draw(dasher, 1)

	var pngBytes bytes.Buffer
	if err := png.Encode(&pngBytes, canvas); err != nil {
		return nil, err
	}

	return pngBytes.Bytes(), nil
}

// openBrowser opens the given URL in the system's default web browser.
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
