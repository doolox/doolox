// Command rendericons rasterizes the Doolox SVG logo into PNGs at several
// sizes so the icon can be installed into the hicolor icon theme and the
// pixmaps directory (both are used by desktop menus such as Cinnamon's).
//
// It reuses the same oksvg/rasterx rendering pipeline that the app itself uses
// for its system-tray icon, so it needs no ImageMagick/rsvg/inkscape.
//
// Usage: go run ./tools/rendericons <svg-path> <output-directory>
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

var sizes = []int{16, 24, 32, 48, 64, 128, 256}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: rendericons <svg-path> <output-directory>")
		os.Exit(2)
	}
	svgPath, outDir := os.Args[1], os.Args[2]

	svgData, err := os.ReadFile(svgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read svg:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir:", err)
		os.Exit(1)
	}

	for _, size := range sizes {
		img, err := render(svgData, size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render %dpx: %v\n", size, err)
			os.Exit(1)
		}
		path := filepath.Join(outDir, fmt.Sprintf("doolox-%d.png", size))
		if err := writePNG(path, img); err != nil {
			fmt.Fprintln(os.Stderr, "write:", err)
			os.Exit(1)
		}
		fmt.Println("wrote", path)
	}
}

// render rasterizes the SVG into a size x size RGBA image.
func render(svgData []byte, size int) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData), oksvg.StrictErrorMode)
	if err != nil {
		return nil, err
	}

	canvas := image.NewRGBA(image.Rect(0, 0, size, size))
	icon.SetTarget(0, 0, float64(size), float64(size))

	scanner := rasterx.NewScannerGV(size, size, canvas, canvas.Bounds())
	dasher := rasterx.NewDasher(size, size, scanner)
	icon.Draw(dasher, 1)

	return canvas, nil
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
