package cms

import (
	"html/template"
	"strings"
	"sync"

	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"

	_ "github.com/go-macaron/session/redis"
)

// renderOpts holds the template configuration shared by the HTTP middleware,
// the system tray menu and the per-website renderers.
var renderOpts macaron.RenderOptions

// siteRenderers caches a compiled renderer per website directory, so each
// website's templates are parsed only once.
var siteRenderers sync.Map

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
		Directory:  contentDir(Website()),
		Extensions: []string{".tmpl", ".html"},
	}

	m.Use(macaron.Renderer(renderOpts))
	m.Use(cache.Cacher())

	return m
}

// rendererWithOpts builds a macaron.Render for the given options.
func rendererWithOpts(opts *macaron.RenderOptions) macaron.Render {
	ts := macaron.NewTemplateSet()
	ts.Set(macaron.DEFAULT_TPL_SET_NAME, opts)

	return &macaron.TplRender{
		TemplateSet:     ts,
		Opt:             opts,
		CompiledCharset: macaron.PrepareCharset(opts.Charset),
	}
}

// newRenderer returns a renderer that uses the current global render options.
// It is used by tests to render a template directory directly.
func newRenderer() macaron.Render {
	return rendererWithOpts(&renderOpts)
}

// rendererFor returns a renderer for the named website. Renderers are cached
// by their template directory so each website's templates are parsed once.
func rendererFor(site string) macaron.Render {
	dir := contentDir(site)
	if r, ok := siteRenderers.Load(dir); ok {
		return r.(macaron.Render)
	}

	opts := renderOpts
	opts.Directory = dir
	r := rendererWithOpts(&opts)
	siteRenderers.Store(dir, r)
	return r
}
