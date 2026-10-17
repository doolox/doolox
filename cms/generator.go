package cms

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	macaron "gopkg.in/macaron.v1"
)

// generateSite generates the static website from the pages in the database.
// It is called from the system tray menu.
func generateSite() {
	GenerateSite()
}

// GenerateSite generates the static website from the pages in the database,
// shows a system notification, and returns the number of generated pages or an
// error.
func GenerateSite() (int, error) {
	generated, err := generate(newRenderer())
	if err != nil {
		log.Printf("[generate] error: %v", err)
		notify("Doolox", fmt.Sprintf("Site generation failed: %v", err))
		return generated, err
	}

	log.Printf("generated %d pages", generated)
	notify("Doolox", fmt.Sprintf("Site generated: %d pages", generated))
	return generated, nil
}

// generate renders every page stored in the database to a static HTML file
// in the website root and returns how many files were written plus the first error
// encountered, if any. The sitemap of the generated pages is written as well.
//
// It does not crawl anything; pages are read from the database. Before writing
// output, generated files are cleared while source and per-site config remain.
func generate(render macaron.Render) (int, error) {
	pages := pages()

	type renderedPage struct {
		page *Page
		html []byte
	}
	rendered := make([]renderedPage, 0, len(pages))
	var firstErr error

	for _, page := range pages {
		html, err := render.HTMLString(
			page.URL,
			map[string]interface{}{"Page": page},
			macaron.HTMLOptions{Layout: layoutFor(page.URL)})
		if err != nil {
			log.Printf("[generate] error while rendering %q: %v", page.URL, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		rendered = append(rendered, renderedPage{page: page, html: []byte(html)})
	}

	if err := clearWebsiteOutput(); err != nil {
		return 0, err
	}

	generated := make([]*Page, 0, len(rendered))
	for _, page := range rendered {
		path := filepath.Join(websiteDir(), page.page.URL+".html")
		if err := os.WriteFile(path, page.html, 0644); err != nil {
			log.Printf("[generate] error while writing %q: %v", path, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		log.Printf("[generate] wrote %q", path)
		generated = append(generated, page.page)
	}

	if err := generateSitemap(generated); err != nil {
		log.Printf("[generate] error while writing the sitemap: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := copyWebsiteAssets(); err != nil {
		log.Printf("[generate] error while copying website assets: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}

	return len(generated), firstErr
}

func clearWebsiteOutput() error {
	if err := os.MkdirAll(websiteDir(), 0755); err != nil {
		return err
	}

	keep := map[string]struct{}{
		"src":         {},
		"config.json": {},
	}
	activeSource := contentDir()
	if filepath.Dir(activeSource) == websiteDir() {
		keep[filepath.Base(activeSource)] = struct{}{}
	}

	entries, err := os.ReadDir(websiteDir())
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if _, preserve := keep[entry.Name()]; preserve {
			continue
		}
		name := strings.ToLower(entry.Name())
		if name == "cname" || strings.TrimSuffix(name, filepath.Ext(name)) == "readme" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(websiteDir(), entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func copyWebsiteAssets() error {
	source := assetsDir()
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	destination := filepath.Join(websiteDir(), "assets")
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported asset type: %s", path)
		}

		input, err := os.Open(path)
		if err != nil {
			return err
		}
		output, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
		if err != nil {
			input.Close()
			return err
		}
		_, copyErr := io.Copy(output, input)
		inputErr := input.Close()
		outputErr := output.Close()
		if copyErr != nil {
			return copyErr
		}
		if inputErr != nil {
			return inputErr
		}
		return outputErr
	})
}

// sitemap is the sitemap.xml document, see https://www.sitemaps.org/protocol.html
type sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	XHTML   string       `xml:"xmlns:xhtml,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// sitemapURL is a single <url> entry of the sitemap.
type sitemapURL struct {
	Loc      string `xml:"loc"`
	Lastmod  string `xml:"lastmod"`
	Priority string `xml:"priority"`
}

// generateSitemap writes sitemap.xml in the website root.
func generateSitemap(pages []*Page) error {
	site := strings.TrimSuffix(conf.URL, "/")

	if site == "" {
		return fmt.Errorf("no url set in the configuration")
	}

	sm := sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		XHTML: "http://www.w3.org/1999/xhtml",
	}

	for _, page := range pages {
		loc := site + "/" + page.URL
		priority := "0.8000"

		if page.URL == indexPage {
			loc = site + "/"
			priority = "1.0000"
		}

		sm.URLs = append(sm.URLs, sitemapURL{
			Loc:      loc,
			Lastmod:  page.Lastmod.UTC().Format("2006-01-02T15:04:05-07:00"),
			Priority: priority,
		})
	}

	body, err := xml.MarshalIndent(sm, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(websiteDir(), "sitemap.xml")
	out := xml.Header + string(body) + "\n"

	if err := os.WriteFile(path, []byte(out), 0644); err != nil {
		return err
	}

	log.Printf("[generate] wrote %q", path)

	return nil
}
