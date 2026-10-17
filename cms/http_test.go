package cms

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/macaron.v1"
)

func TestSimpleSeedLayoutsRender(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate test source")
	}
	content := filepath.Join(filepath.Dir(sourceFile), "..", "defaults", "simple", "src")

	previousOpts := renderOpts
	renderOpts.Directory = content
	defer func() { renderOpts = previousOpts }()

	renderer := newRenderer()
	for _, pageURL := range []string{indexPage, "about"} {
		t.Run(pageURL, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(content, "pages", pageURL+".html"))
			if err != nil {
				t.Fatal(err)
			}
			page := newPage(pageURL, data)
			html, err := renderer.HTMLString(pageURL, map[string]interface{}{"Page": page}, macaron.HTMLOptions{Layout: layoutFor(pageURL)})
			if err != nil {
				t.Fatalf("render %s layout: %v", layoutFor(pageURL), err)
			}
			if !strings.Contains(html, "/assets/vendor/bootstrap/css/bootstrap.min.css") {
				t.Errorf("rendered %s page is missing Bootstrap CSS", pageURL)
			}
			if !strings.Contains(html, "/assets/img/favicon.ico") {
				t.Errorf("rendered %s page is missing the Doolox favicon", pageURL)
			}
			if !strings.Contains(html, "/assets/img/favicon.svg") {
				t.Errorf("rendered %s page is missing the Simple SVG favicon", pageURL)
			}
			if pageURL == indexPage && !strings.Contains(html, "framework: bootstrap") {
				t.Error("rendered homepage is missing the software preview")
			}
		})
	}
}

func TestDooloxSeedLayoutsRender(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate test source")
	}
	content := filepath.Join(filepath.Dir(sourceFile), "..", "defaults", "doolox", "src")

	previousOpts := renderOpts
	renderOpts.Directory = content
	defer func() { renderOpts = previousOpts }()

	renderer := newRenderer()
	for _, pageURL := range []string{indexPage, "about"} {
		t.Run(pageURL, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(content, "pages", pageURL+".html"))
			if err != nil {
				t.Fatal(err)
			}
			page := newPage(pageURL, data)
			html, err := renderer.HTMLString(pageURL, map[string]interface{}{"Page": page}, macaron.HTMLOptions{Layout: layoutFor(pageURL)})
			if err != nil {
				t.Fatalf("render %s layout: %v", layoutFor(pageURL), err)
			}
			if !strings.Contains(html, "/assets/img/favicon.ico") {
				t.Errorf("rendered %s page is missing the Doolox favicon", pageURL)
			}
			if !strings.Contains(html, "assets/vendor/bootstrap/css/bootstrap.min.css") {
				t.Errorf("rendered %s page is missing Bootstrap CSS", pageURL)
			}
		})
	}
}

func TestContentDirPrefersSourceAndFallsBackToLegacy(t *testing.T) {
	t.Setenv(dooloxHomeEnv, t.TempDir())
	previousConf := conf
	conf = &Config{Website: "simple"}
	defer func() { conf = previousConf }()

	source := sourceDir()
	if err := os.MkdirAll(filepath.Join(source, "assets"), 0755); err != nil {
		t.Fatal(err)
	}
	if got := contentDir(); got != source {
		t.Fatalf("contentDir() = %q, want source directory %q", got, source)
	}
	if got, want := assetsDir(), filepath.Join(source, "assets"); got != want {
		t.Fatalf("assetsDir() = %q, want %q", got, want)
	}

	if err := os.RemoveAll(source); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(websiteDir(), "content")
	if err := os.MkdirAll(legacy, 0755); err != nil {
		t.Fatal(err)
	}
	if got := contentDir(); got != legacy {
		t.Fatalf("contentDir() = %q, want legacy directory %q", got, legacy)
	}
	if got, want := assetsDir(), filepath.Join(legacy, "theme", "assets"); got != want {
		t.Fatalf("assetsDir() = %q, want legacy assets directory %q", got, want)
	}
}

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
