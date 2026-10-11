package cms

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestRenderSVGIcon(t *testing.T) {
	svg, err := os.ReadFile("../assets/trayicon.svg")
	if err != nil {
		t.Fatalf("read tray icon SVG: %v", err)
	}

	data, err := renderSVGIcon(svg)
	if err != nil {
		t.Fatalf("render tray icon: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode rendered tray icon: %v", err)
	}
	if got, want := img.Bounds(), image.Rect(0, 0, 64, 64); got != want {
		t.Fatalf("rendered bounds = %v, want %v", got, want)
	}
	if _, _, _, alpha := img.At(0, 0).RGBA(); alpha != 0 {
		t.Fatalf("corner alpha = %d, want transparent", alpha)
	}
	r, g, b, alpha := img.At(32, 16).RGBA()
	if alpha == 0 || r <= g || g <= b {
		t.Fatalf("top face pixel = (%d, %d, %d, %d), want opaque yellow", r, g, b, alpha)
	}
}
