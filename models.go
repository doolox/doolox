package main

import (
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	// PAGE type of page
	PAGE = "page"
	// POST type of page
	POST = "post"
)

// Page model represents any kind of page
type Page struct {
	gorm.Model
	URL         string `sql:"size:255;unique_index"`
	Title       string `sql:"size:255"`
	Description string `sql:"size:255"`
	Type        string `sql:"size:255"`
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
