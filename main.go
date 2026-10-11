package main

import (
	_ "embed"

	"doolox/cms"
)

//go:embed assets/trayicon.png
var trayIconPNG []byte

//go:embed assets/trayicon.ico
var trayIconICO []byte

func main() {
	cms.Start(cms.TrayIcon{
		PNG: trayIconPNG,
		ICO: trayIconICO,
	})
}
