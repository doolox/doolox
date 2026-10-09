package cms

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	macaron "gopkg.in/macaron.v1"
)

// docsDir is the folder the static website is generated into.
const docsDir = "docs"

// generate renders every page stored in the database to a static HTML file
// inside docs and returns how many files were written.
//
// It does not crawl anything, the pages are read from the database. Files that
// already exist are overridden, folders are never created or copied, only
// plain HTML files are written.
func generate(render macaron.Render) int {
	pages := []*Page{}
	db.Find(&pages)

	generated := 0

	for _, page := range pages {
		if strings.Contains(page.URL, "/") {
			log.Printf("[generate] %q skipped, only plain html files are generated", page.URL)
			continue
		}

		html, err := render.HTMLString(
			page.URL,
			map[string]interface{}{"Page": page},
			macaron.HTMLOptions{Layout: layoutFor(page.URL)})
		if err != nil {
			log.Printf("[generate] error while rendering %q: %v", page.URL, err)
			continue
		}

		path := filepath.Join(docsDir, page.URL+".html")

		if err := os.WriteFile(path, []byte(html), 0644); err != nil {
			log.Printf("[generate] error while writing %q: %v", path, err)
			continue
		}

		log.Printf("[generate] wrote %q", path)
		generated++
	}

	return generated
}
