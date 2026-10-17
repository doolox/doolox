package main

import (
	_ "embed"
	"flag"

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
	openWebsite := flag.Bool("open-website", true, "open the Doolox website on startup")
	flag.Parse()

	go admin.Start()

	cms.Start(cms.TrayIcon{
		SVG: trayIconSVG,
		PNG: trayIconPNG,
		ICO: trayIconICO,
	}, *openWebsite)
}
