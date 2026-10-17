package cms

import (
	"fmt"
	"log"
	"net/http"

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

func pageView(ctx *macaron.Context) {
	page := ctx.Params(":page")
	if page == "" {
		page = indexPage
	}

	p, ok := pageByURL(page)
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.Data["Page"] = p

	// Render to a string so the Doolox preview helper can be injected for
	// live serving only (the static site generator renders without it).
	html, err := ctx.Render.HTMLString(page, ctx.Data, macaron.HTMLOptions{Layout: layoutFor(page)})
	if err != nil {
		log.Printf("[pageView] render %q: %v", page, err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Resp.WriteHeader(http.StatusOK)
	ctx.Resp.Write([]byte(injectDooloxJS(html)))
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
