package cms

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHTTPRoutes(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(dooloxHomeEnv, dir)

	prevConf := conf
	prevM := m
	prevOpts := renderOpts
	conf = &Config{Website: "doolox", URL: "https://example.com"}
	defer func() {
		conf = prevConf
		m = prevM
		renderOpts = prevOpts
	}()

	content := filepath.Join(dir, "websites", "doolox", "content")
	mk := func(rel, body string) {
		t.Helper()
		p := filepath.Join(content, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	mk("pages/index.html", "<!--\ntitle: Doolox\n-->\n{{define \"index\"}}{{end}}")
	mk("pages/about.html", "<!--\ntitle: About Doolox\nsubtitle: About\n-->\n{{define \"about\"}}ABOUT BODY{{end}}")
	mk("layout_home.html", "<html><head>{{template \"header\" .}}</head><body>HOME</body></html>")
	mk("layout.html", "<html><head>{{template \"header\" .}}</head><body>{{yield}}</body></html>")
	mk("theme/header.html", "{{define \"header\"}}<title>{{.Page.Title}}</title>{{end}}")

	m = initMacaron()
	m.Get("/", pageView)
	m.Get("/:page", pageView)

	get := func(path string) *httptest.ResponseRecorder {
		t.Helper()
		rec := httptest.NewRecorder()
		m.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}

	rec := get("/")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Doolox") {
		t.Errorf("GET / body missing title %q", rec.Body.String())
	}

	rec = get("/about")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /about = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "About Doolox") {
		t.Errorf("GET /about body missing title: %q", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "ABOUT BODY") {
		t.Errorf("GET /about body missing page content: %q", rec.Body.String())
	}
}

func TestHTTPRouteLegacyHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(dooloxHomeEnv, dir)

	prevConf := conf
	prevM := m
	prevOpts := renderOpts
	conf = &Config{Website: "doolox", URL: "https://example.com"}
	defer func() {
		conf = prevConf
		m = prevM
		renderOpts = prevOpts
	}()

	content := filepath.Join(dir, "websites", "doolox", "content")
	mk := func(rel, body string) {
		t.Helper()
		p := filepath.Join(content, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// Legacy home.html (old filename) with no front matter.
	mk("pages/home.html", "")
	mk("layout_home.html", "<html><head><title>HOME</title></head><body>HOME</body></html>")

	m = initMacaron()
	m.Get("/", pageView)
	m.Get("/:page", pageView)

	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "HOME") {
		t.Errorf("GET / body missing home content: %q", rec.Body.String())
	}
}
