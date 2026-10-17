package admin

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"doolox/cms"
)

// addr is the address the admin web server listens on.
const addr = ":5555"

// staticDir returns the directory holding the static assets (css, js, images,
// ...) served under /assets. The Windows installer and the Debian package ship
// these files next to the executable (admin/theme/static), so for an installed
// binary the directory is resolved relative to the executable's own location.
// When running from the source tree (`go run .`), it falls back to the
// repository-relative path.
func staticDir() string {
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Join(filepath.Dir(exe), "admin", "theme", "static")
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}
	return filepath.Join("admin", "theme", "static")
}

// versionFile is the single source of truth for the Doolox version. It lives
// next to the repository root and is read at runtime relative to the working
// directory.
const versionFile = "version.iss"

//go:embed theme/templates/*.html
var templatesFS embed.FS

//go:embed theme/static/logo.svg
var logo []byte

// indexTmpl is the AdminLTE layout served at "/".
var indexTmpl = template.Must(template.ParseFS(templatesFS, "theme/templates/*.html"))

// pageData holds the values injected into the layout template.
type pageData struct {
	Version string
	Pages   []cms.Page
	Message string
	Error   bool
}

// appVersion reads the Doolox version from version.iss. It returns "dev" when
// the file is missing or does not contain a version.
func appVersion() string {
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return "dev"
	}
	if v := parseVersion(data); v != "" {
		return v
	}
	return "dev"
}

// parseVersion extracts the version from version.iss content, e.g.
// `#define MyAppVersion "1.0.1"` -> "1.0.1".
func parseVersion(data []byte) string {
	const prefix = `#define MyAppVersion "`

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, `"`) {
			return strings.TrimSuffix(strings.TrimPrefix(line, prefix), `"`)
		}
	}

	return ""
}

// Start runs the admin web server on port 5555. It blocks until the server
// stops, so run it in its own goroutine when other servers must run too.
func Start() {
	server := &http.Server{Addr: addr, Handler: handler()}

	log.Printf("[admin] listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("[admin] server error: %v", err)
	}
}

// handler builds the admin HTTP handler: the AdminLTE static assets, the
// Doolox logo, and the single index page.
func handler() http.Handler {
	mux := http.NewServeMux()

	// Static assets (AdminLTE css/js, ...) served from admin/theme/static.
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir()))))

	// Doolox logo (embedded in the binary).
	mux.HandleFunc("/logo.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(logo)
	})

	// Delete a page (POST) and redirect back to the index.
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		url := r.FormValue("url")
		if url == "" {
			http.Error(w, "missing url", http.StatusBadRequest)
			return
		}
		if err := cms.DeletePage(url); err != nil {
			log.Printf("[admin] delete %q: %v", url, err)
			http.Error(w, "could not delete page", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Generate the static website and redirect back to the index with a success
	// or error message.
	mux.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		generated, err := cms.GenerateSite()
		if err != nil {
			http.Redirect(w, r, "/?status=error&msg="+url.QueryEscape(err.Error()), http.StatusSeeOther)
			return
		}
		msg := fmt.Sprintf("Site generated: %d pages", generated)
		http.Redirect(w, r, "/?status=success&msg="+url.QueryEscape(msg), http.StatusSeeOther)
	})

	// The admin index page.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := pageData{
			Version: appVersion(),
			Pages:   cms.Pages(),
			Message: r.URL.Query().Get("msg"),
			Error:   r.URL.Query().Get("status") == "error",
		}
		if err := indexTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
			log.Printf("[admin] render index: %v", err)
		}
	})

	return mux
}
