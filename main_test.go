package main

import (
	"bytes"
	"testing"
)

func TestTrayIconSVGEmbedded(t *testing.T) {
	if !bytes.Contains(trayIconSVG, []byte("<svg")) {
		t.Fatal("tray SVG was not embedded")
	}
}
