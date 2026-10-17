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
	generated, err := generate()
	if err != nil {
		log.Printf("[generate] error: %v", err)
		notify("Doolox", fmt.Sprintf("Site generation failed: %v", err))
		return generated, err
	}

	log.Printf("generated %d pages", generated)
	notify("Doolox", fmt.Sprintf("Site generated: %d pages", generated))
	return generated, nil
}

// sitePage pairs a page with the website it belongs to, for the combined
// sitemap.
type sitePage struct {
	site string
	page *Page
}

// generate renders every website to a static HTML file inside its own folder
// and writes a single combined sitemap.xml at the root of the websites
// directory. It returns how many files were written plus the first error
// encountered, if any.
//
// It does not crawl anything; pages are read from disk. Before writing output,
// each website's generated files are cleared while its source and per-site
// config remain.
func generate() (int, error) {
	total := 0
	allPages := make([]sitePage, 0)
	var firstErr error

	for _, site := range websites() {
		generated, sitePages, err := generateWebsite(site)
		total += generated
		for _, p := range sitePages {
			allPages = append(allPages, sitePage{site: site, page: p})
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if err := generateSitemap(allPages); err != nil {
		log.Printf("[generate] error while writing the sitemap: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}

	return total, firstErr
}

// generateWebsite renders the pages of a single website to static HTML files
// inside its folder and copies its assets, returning the generated pages.
func generateWebsite(site string) (int, []*Page, error) {
	renderer := rendererFor(site)
	sitePages := pages(site)

	if err := clearSiteOutput(site); err != nil {
		return 0, nil, err
	}

	generated := make([]*Page, 0, len(sitePages))
	var firstErr error
	for _, page := range sitePages {
		html, err := renderer.HTMLString(
			page.URL,
			map[string]interface{}{"Page": page},
			macaron.HTMLOptions{Layout: layoutFor(page.URL)})
		if err != nil {
			log.Printf("[generate] error while rendering %s/%s: %v", site, page.URL, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		path := filepath.Join(websiteDir(site), page.URL+".html")
		if err := os.WriteFile(path, []byte(html), 0644); err != nil {
			log.Printf("[generate] error while writing %q: %v", path, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		log.Printf("[generate] wrote %q", path)
		generated = append(generated, page)
	}

	if err := copySiteAssets(site); err != nil {
		log.Printf("[generate] error while copying website assets: %v", err)
		if firstErr == nil {
			firstErr = err
		}
	}

	return len(generated), generated, firstErr
}

func clearSiteOutput(site string) error {
	dir := websiteDir(site)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	keep := map[string]struct{}{
		"src":         {},
		"config.json": {},
		".git":        {},
	}
	activeSource := contentDir(site)
	if filepath.Dir(activeSource) == dir {
		keep[filepath.Base(activeSource)] = struct{}{}
	}
	// Nested websites are subfolders of the default website and must survive
	// the default website's output clearing.
	for _, n := range websites() {
		keep[n] = struct{}{}
	}

	entries, err := os.ReadDir(dir)
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
		if err := os.RemoveAll(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func copySiteAssets(site string) error {
	source := assetsDir(site)
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	destination := filepath.Join(websiteDir(site), "assets")
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

// generateSitemap writes a single combined sitemap.xml at the root of the
// websites directory, listing every page of every website under its own path.
func generateSitemap(sitePages []sitePage) error {
	site := strings.TrimSuffix(conf.URL, "/")

	if site == "" {
		return fmt.Errorf("no url set in the configuration")
	}

	sm := sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		XHTML: "http://www.w3.org/1999/xhtml",
	}

	for _, sp := range sitePages {
		priority := "0.8000"
		if sp.page.URL == indexPage {
			priority = "1.0000"
		}

		// The default website is served at the root of the URL; every other
		// website is served under its own subpath.
		var loc string
		if sp.site == Website() {
			loc = site + "/"
		} else {
			loc = site + "/" + sp.site + "/"
		}
		if sp.page.URL != indexPage {
			loc += sp.page.URL
		}

		sm.URLs = append(sm.URLs, sitemapURL{
			Loc:      loc,
			Lastmod:  sp.page.Lastmod.UTC().Format("2006-01-02T15:04:05-07:00"),
			Priority: priority,
		})
	}

	body, err := xml.MarshalIndent(sm, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(websitesDir(), "sitemap.xml")
	out := xml.Header + string(body) + "\n"

	if err := os.WriteFile(path, []byte(out), 0644); err != nil {
		return err
	}

	log.Printf("[generate] wrote %q", path)

	return nil
}
