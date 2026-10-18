package cms

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	macaron "gopkg.in/macaron.v1"
)

// indexPage is the url of the page served for "/".
const indexPage = "index"

// layoutFor returns the layout used to render the page with the given url.
func layoutFor(url string) string {
	if url == indexPage {
		return "layout_home"
	}

	return "layout"
}

// homeView renders the home page of the default website.
func homeView(ctx *macaron.Context) {
	renderPage(ctx, Website(), indexPage)
}

// segmentView handles a single path segment: when it names a website other
// than the default one it renders that website's home page, otherwise the
// segment is treated as a page of the default website.
func segmentView(ctx *macaron.Context) {
	seg := ctx.Params("seg")
	if isWebsite(seg) && seg != Website() {
		renderPage(ctx, seg, indexPage)
		return
	}
	renderPage(ctx, Website(), seg)
}

// nestedView renders a page of the website named by the first path segment.
func nestedView(ctx *macaron.Context) {
	seg := ctx.Params("seg")
	page := ctx.Params("page")
	if !isWebsite(seg) || seg == Website() {
		ctx.Status(http.StatusNotFound)
		return
	}
	renderPage(ctx, seg, page)
}

// renderPage renders the given page of the given website.
func renderPage(ctx *macaron.Context, site, page string) {
	p, ok := pageByURL(site, page)
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.Data["Page"] = p

	// Render to a string so the Doolox preview helper can be injected for
	// live serving only (the static site generator renders without it).
	html, err := rendererFor(site).HTMLString(page, ctx.Data, macaron.HTMLOptions{Layout: layoutFor(page)})
	if err != nil {
		log.Printf("[renderPage] render %s/%s: %v", site, page, err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Resp.WriteHeader(http.StatusOK)
	ctx.Resp.Write([]byte(injectDooloxJS(html)))
}

// assetsView serves the static assets of the default website.
func assetsView(ctx *macaron.Context) {
	serveAssets(ctx, Website())
}

// nestedAssetsView serves the static assets of the website named by the first
// path segment.
func nestedAssetsView(ctx *macaron.Context) {
	site := ctx.Params("seg")
	if !isWebsite(site) {
		ctx.Status(http.StatusNotFound)
		return
	}
	serveAssets(ctx, site)
}

// serveAssets serves a file from the given website's assets directory.
func serveAssets(ctx *macaron.Context, site string) {
	dir := filepath.Clean(assetsDir(site))
	name := filepath.Join(dir, filepath.FromSlash(ctx.Params("*")))

	// Keep the request inside the assets directory.
	if name != dir && !strings.HasPrefix(name, dir+string(os.PathSeparator)) {
		ctx.Status(http.StatusNotFound)
		return
	}

	http.ServeFile(ctx.Resp, ctx.Req.Request, name)
}

func generateView(ctx *macaron.Context) {
	// GenerateSite also shows a native system notification (the same one the
	// tray menu's "Generate" action uses).
	generated, err := GenerateSite()
	if err != nil {
		ctx.Render.PlainText(
			http.StatusInternalServerError,
			[]byte(fmt.Sprintf("generation failed: %v\n", err)))
		return
	}

	ctx.Render.PlainText(
		http.StatusOK,
		[]byte(fmt.Sprintf("generated %d pages\n", generated)))
}

// publishView generates the site and then commits and pushes it with git. It
// is the endpoint behind the Ctrl+P preview shortcut.
func publishView(ctx *macaron.Context) {
	msg, err := PublishSite()
	if err != nil {
		ctx.Render.PlainText(
			http.StatusInternalServerError,
			[]byte(fmt.Sprintf("publish failed: %v\n", err)))
		return
	}

	ctx.Render.PlainText(
		http.StatusOK,
		[]byte(msg+"\n"))
}

// dashboardView opens the Doolox Dashboard the same way the tray menu's
// "Doolox Dashboard" action does: a native web view on Windows and macOS, the
// system browser on Linux. It is the endpoint behind the Ctrl+D preview
// shortcut. launchDashboard blocks on the web view platforms, so run it in its
// own goroutine.
func dashboardView(ctx *macaron.Context) {
	go launchDashboard()
	ctx.Render.PlainText(http.StatusOK, []byte("opening dashboard\n"))
}
