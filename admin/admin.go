package admin

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// addr is the address the admin web server listens on.
const addr = ":5555"

// staticDir is the directory holding the static assets (css, js, images, ...)
// served under /assets. These files are committed and tracked by git.
const staticDir = "admin/theme/static"

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
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticDir))))

	// Doolox logo (embedded in the binary).
	mux.HandleFunc("/logo.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(logo)
	})

	// The single admin page for now; everything else is a 404.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTmpl.ExecuteTemplate(w, "layout.html", pageData{Version: appVersion()}); err != nil {
			log.Printf("[admin] render index: %v", err)
		}
	})

	return mux
}
