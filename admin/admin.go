package admin

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

// addr is the address the admin web server listens on.
const addr = ":5555"

// adminlteDist is the directory holding the AdminLTE distribution (css, js,
// images, ...). It is ignored by git and populated by unpacking the AdminLTE
// release archive into admin/theme/static.
const adminlteDist = "admin/theme/static/adminlte"

//go:embed theme/templates/*.html
var templatesFS embed.FS

//go:embed theme/static/logo.svg
var logo []byte

// indexTmpl is the AdminLTE layout served at "/".
var indexTmpl = template.Must(template.ParseFS(templatesFS, "theme/templates/*.html"))

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

	// AdminLTE static assets served from admin/theme/static/adminlte.
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(adminlteDist))))

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
		if err := indexTmpl.ExecuteTemplate(w, "layout.html", nil); err != nil {
			log.Printf("[admin] render index: %v", err)
		}
	})

	return mux
}
