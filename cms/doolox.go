package cms

import (
	_ "embed"
	"net/http"
	"strings"

	macaron "gopkg.in/macaron.v1"
)

//go:embed doolox.js
var dooloxJS []byte

// dooloxJSTag is the script tag injected into pages served live by Doolox so
// the preview helper (the Ctrl+G shortcut that generates the site) is
// available.
const dooloxJSTag = `<script src="/doolox.js"></script>`

// injectDooloxJS inserts the Doolox helper script into a rendered page just
// before </body>. It is only used for live preview; the static site generator
// never injects it.
func injectDooloxJS(html string) string {
	if i := strings.LastIndex(html, "</body>"); i >= 0 {
		return html[:i] + dooloxJSTag + html[i:]
	}
	return html + dooloxJSTag
}

// dooloxJSView serves the Doolox live-preview helper script.
func dooloxJSView(ctx *macaron.Context) {
	ctx.Resp.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	ctx.Resp.WriteHeader(http.StatusOK)
	ctx.Resp.Write(dooloxJS)
}
