package cms

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// pagesDir is the directory holding the website pages. Each .html file in it
// is one page; the page URL is the file name without the extension.
func pagesDir() string { return filepath.Join(contentDir(), "pages") }

// parseFrontMatter extracts the key/value pairs from the optional leading
// <!-- ... --> comment of a page file. It returns nil when the file has no
// front matter.
func parseFrontMatter(data []byte) map[string]string {
	text := strings.TrimSpace(string(data))
	if !strings.HasPrefix(text, "<!--") {
		return nil
	}
	end := strings.Index(text, "-->")
	if end < 0 {
		return nil
	}

	meta := map[string]string{}
	for _, line := range strings.Split(text[len("<!--"):end], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, ":")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			continue
		}
		meta[strings.ToLower(key)] = strings.TrimSpace(val)
	}
	return meta
}

// newPage builds a Page for the given url from the raw page file content,
// filling in the metadata from the front matter.
func newPage(url string, data []byte) *Page {
	p := &Page{URL: url}
	for key, val := range parseFrontMatter(data) {
		switch key {
		case "title":
			p.Title = val
		case "subtitle":
			p.Subtitle = val
		case "description":
			p.Description = val
		case "type":
			p.Type = val
		}
	}
	return p
}

// pages returns every page discovered in the pages directory, ordered by URL.
func pages() []*Page {
	dir := pagesDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("[pages] read %s: %v", dir, err)
		return nil
	}

	// Keyed by URL so that when both a legacy home.html and an index.html
	// exist, the canonical index.html wins: entries are returned in filename
	// order and "index.html" sorts after "home.html", so it overwrites it.
	byURL := map[string]*Page{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".html" {
			continue
		}
		name := entry.Name()
		url := strings.TrimSuffix(name, filepath.Ext(name))
		// A legacy home.html represents the same page as index.html (the root).
		if url == "home" {
			url = indexPage
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			log.Printf("[pages] read %s: %v", name, err)
			continue
		}

		page := newPage(url, data)
		if info, err := entry.Info(); err == nil {
			page.Lastmod = info.ModTime()
		}
		byURL[url] = page
	}

	list := make([]*Page, 0, len(byURL))
	for _, p := range byURL {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].URL < list[j].URL })
	return list
}

// pageByURL returns the page with the given URL, or false when no such page
// exists.
func pageByURL(url string) (*Page, bool) {
	for _, p := range pages() {
		if p.URL == url {
			return p, true
		}
	}
	return nil, false
}
