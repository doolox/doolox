package cms

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInjectDooloxJS(t *testing.T) {
	in := "<html><head></head><body>hello</body></html>"
	want := "<html><head></head><body>hello" + dooloxJSTag + "</body></html>"
	if got := injectDooloxJS(in); got != want {
		t.Errorf("injectDooloxJS() = %q, want %q", got, want)
	}

	noBody := "<html><head></head></html>"
	if got := injectDooloxJS(noBody); got != noBody+dooloxJSTag {
		t.Errorf("injectDooloxJS() without </body> = %q, want %q", got, noBody+dooloxJSTag)
	}
}

func TestDooloxJSServedAndInjected(t *testing.T) {
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
	mk("layout_home.html", "<html><head>{{template \"header\" .}}</head><body>HOME</body></html>")
	mk("theme/header.html", "{{define \"header\"}}<title>{{.Page.Title}}</title>{{end}}")

	m = initMacaron()
	m.Get("/", pageView)
	m.Get("/doolox.js", dooloxJSView)
	m.Get("/:page", pageView)

	// The live page must include the helper script tag.
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), dooloxJSTag) {
		t.Errorf("GET / body missing doolox script tag: %q", rec.Body.String())
	}

	// The script must be served with the correct content type.
	rec2 := httptest.NewRecorder()
	m.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/doolox.js", nil))
	if rec2.Code != http.StatusOK {
		t.Fatalf("GET /doolox.js = %d, want 200", rec2.Code)
	}
	if ct := rec2.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/javascript") {
		t.Errorf("Content-Type = %q, want application/javascript", ct)
	}
	body := rec2.Body.String()
	if !strings.Contains(body, "Ctrl+G") {
		t.Error("/doolox.js body missing Ctrl+G handling")
	}
	if !strings.Contains(body, "/generate") {
		t.Error("/doolox.js body missing /generate endpoint")
	}
}
