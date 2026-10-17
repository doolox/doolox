package main

import (
	_ "embed"

	"doolox/admin"
	"doolox/cms"
)

//go:embed assets/trayicon.svg
var trayIconSVG []byte

//go:embed assets/trayicon.png
var trayIconPNG []byte

//go:embed assets/trayicon.ico
var trayIconICO []byte

func main() {
	go admin.Start()

	cms.Start(cms.TrayIcon{
		SVG: trayIconSVG,
		PNG: trayIconPNG,
		ICO: trayIconICO,
	})
}
