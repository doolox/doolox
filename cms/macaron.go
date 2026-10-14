package cms

import (
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"

	_ "github.com/go-macaron/session/redis"
)

// renderOpts holds the template configuration shared by the HTTP middleware
// and the system tray menu.
var renderOpts macaron.RenderOptions

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	renderOpts = macaron.RenderOptions{
		Layout: "layout",
		Funcs: []template.FuncMap{map[string]interface{}{
			"obfuscate": func(args ...interface{}) template.HTML {
				email := args[0].(string)
				email = strings.Replace(email, "@", "<span style=\"display:none\">evilspam</span>@", 1)
				return template.HTML(email)

			},
		}},
		Directory:  contentDir(),
		Extensions: []string{".tmpl", ".html"},
	}

	s := macaron.Static(
		filepath.Join(contentDir(), "theme", "assets"),
		macaron.StaticOptions{
			Prefix:      "assets",
			SkipLogging: true,
			IndexFile:   "index.html",
			// Expires defines which user-defined function to use for producing a HTTP Expires Header. Default is nil.
			// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
			Expires: func() string {
				return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
			},
		})

	m.Use(macaron.Renderer(renderOpts))
	m.Use(s)
	m.Use(cache.Cacher())

	return m
}

// newRenderer returns a macaron.Render that can render templates outside of an
// HTTP request. It is used by the system tray menu to generate the site.
func newRenderer() macaron.Render {
	ts := macaron.NewTemplateSet()
	ts.Set(macaron.DEFAULT_TPL_SET_NAME, &renderOpts)

	return &macaron.TplRender{
		TemplateSet:     ts,
		Opt:             &renderOpts,
		CompiledCharset: macaron.PrepareCharset(renderOpts.Charset),
	}
}
