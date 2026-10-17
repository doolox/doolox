package cms

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]string
	}{
		{
			name: "standard block",
			in:   "<!--\ntitle: About Doolox\nsubtitle: About\n-->",
			want: map[string]string{
				"title":    "About Doolox",
				"subtitle": "About",
			},
		},
		{
			name: "whitespace and case",
			in:   "<!--\n  Title:  Hello  \nTYPE: page\n-->",
			want: map[string]string{
				"title": "Hello",
				"type":  "page",
			},
		},
		{
			name: "ignores malformed lines",
			in:   "<!--\ntitle: X\nno colon here\n: empty key\n-->",
			want: map[string]string{"title": "X"},
		},
		{
			name: "no front matter",
			in:   "{{define \"about\"}}hello{{end}}",
			want: nil,
		},
		{
			name: "leading whitespace then comment",
			in:   "\n  <!--\ntitle: X\n-->\n{{define \"x\"}}{{end}}",
			want: map[string]string{"title": "X"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFrontMatter([]byte(tt.in))
			if len(got) != len(tt.want) {
				t.Fatalf("parseFrontMatter() = %v, want %v", got, tt.want)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("parseFrontMatter()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	data := []byte("<!--\ntitle: About Doolox\nsubtitle: About\ndescription: About Doolox\ntype: page\n-->\n{{define \"about\"}}...{{end}}")

	p := newPage("about", data)
	if p.URL != "about" {
		t.Errorf("URL = %q, want %q", p.URL, "about")
	}
	if p.Title != "About Doolox" {
		t.Errorf("Title = %q, want %q", p.Title, "About Doolox")
	}
	if p.Subtitle != "About" {
		t.Errorf("Subtitle = %q, want %q", p.Subtitle, "About")
	}
	if p.Description != "About Doolox" {
		t.Errorf("Description = %q, want %q", p.Description, "About Doolox")
	}
	if p.Type != "page" {
		t.Errorf("Type = %q, want %q", p.Type, "page")
	}
}

func TestGenerateWritesToWebsiteRoot(t *testing.T) {
	t.Setenv(dooloxHomeEnv, t.TempDir())
	previousConf := conf
	previousOpts := renderOpts
	conf = &Config{Website: "simple", URL: "https://example.com"}
	defer func() {
		conf = previousConf
		renderOpts = previousOpts
	}()

	source := sourceDir()
	pagesPath := filepath.Join(source, "pages")
	if err := os.MkdirAll(pagesPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pagesPath, "index.html"), []byte("<!--\ntitle: Home\n-->\n{{define \"index\"}}Home page{{end}}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "layout_home.html"), []byte("<html><body>{{yield}}</body></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	renderOpts.Directory = source

	generated, err := generate(newRenderer())
	if err != nil {
		t.Fatalf("generate() error = %v", err)
	}
	if generated != 1 {
		t.Fatalf("generate() wrote %d pages, want 1", generated)
	}
	for _, name := range []string{"index.html", "sitemap.xml"} {
		if _, err := os.Stat(filepath.Join(websiteDir(), name)); err != nil {
			t.Errorf("generated file %q missing from website root: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(websiteDir(), "docs")); !os.IsNotExist(err) {
		t.Errorf("docs directory exists after generation, stat error = %v", err)
	}
}

func TestPagesLegacyHomeMapsToIndex(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(dooloxHomeEnv, dir)

	prev := conf
	conf = &Config{Website: "doolox"}
	defer func() { conf = prev }()

	pagesDir := filepath.Join(dir, "websites", "doolox", "content", "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pagesDir, "about.html"), []byte("<!--\ntitle: About\n-->\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// A legacy home.html (old filename) should resolve as the index page.
	if err := os.WriteFile(filepath.Join(pagesDir, "home.html"), []byte("<!--\ntitle: Legacy Home\n-->\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ps := pages()
	if len(ps) != 2 {
		t.Fatalf("pages() returned %d pages, want 2: %+v", len(ps), ps)
	}
	byURL := map[string]*Page{}
	for _, p := range ps {
		byURL[p.URL] = p
	}
	if got := byURL["about"].Title; got != "About" {
		t.Errorf("about title = %q, want %q", got, "About")
	}
	idx, ok := byURL[indexPage]
	if !ok {
		t.Fatalf("missing %q page: %+v", indexPage, ps)
	}
	if idx.Title != "Legacy Home" {
		t.Errorf("index title = %q, want %q", idx.Title, "Legacy Home")
	}
}

func TestPagesIndexWinsOverLegacyHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(dooloxHomeEnv, dir)

	prev := conf
	conf = &Config{Website: "doolox"}
	defer func() { conf = prev }()

	pagesDir := filepath.Join(dir, "websites", "doolox", "content", "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pagesDir, "home.html"), []byte("<!--\ntitle: Legacy\n-->\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pagesDir, "index.html"), []byte("<!--\ntitle: Canonical\n-->\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ps := pages()
	if len(ps) != 1 {
		t.Fatalf("pages() returned %d pages, want 1: %+v", len(ps), ps)
	}
	if ps[0].URL != indexPage || ps[0].Title != "Canonical" {
		t.Errorf("pages()[0] = %+v, want index/Canonical", ps[0])
	}
}
