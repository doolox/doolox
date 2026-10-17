package cms

import (
	"strings"
	"time"
)

const (
	// PAGE type of page
	PAGE = "page"
	// POST type of page
	POST = "post"
)

// Page represents a single page of the website. Pages are defined by files in
// the content/pages directory: the URL is the file name without its extension
// and the metadata (title, subtitle, description, type) is read from the
// optional front-matter comment at the top of the file.
type Page struct {
	URL         string
	Title       string
	Subtitle    string
	Description string
	Type        string
	Lastmod     time.Time
}

// IsActive reports whether the given menu link points to this page.
// The root link "/" is mapped to the "index" page.
func (p *Page) IsActive(link string) bool {
	if p == nil {
		return false
	}

	slug := strings.Trim(link, "/")
	if slug == "" {
		slug = "index"
	}

	return p.URL == slug
}
