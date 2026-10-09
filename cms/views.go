package cms

import (
	"fmt"
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

	p := &Page{URL: page}
	db.FirstOrCreate(p, p)
	ctx.Data["Page"] = p

	ctx.HTML(
		200,
		page,
		ctx.Data,
		macaron.HTMLOptions{Layout: layoutFor(page)})
}

func generateView(ctx *macaron.Context) {
	generated := generate(ctx.Render)

	ctx.Render.PlainText(
		http.StatusOK,
		[]byte(fmt.Sprintf("generated %d pages\n", generated)))
}
