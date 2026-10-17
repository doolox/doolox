package admin

import (
	"net/http"
	"net/url"
	"strings"
)

// flashCookie is the cookie used to carry a one-shot flash message from an
// action (such as /generate) to the index page across the redirect that
// follows it.
const flashCookie = "doolox_flash"

// setFlash stores a flash message in a short-lived cookie so it survives the
// redirect to the index page and is shown exactly once.
func setFlash(w http.ResponseWriter, msg string, isError bool) {
	kind := "success"
	if isError {
		kind = "error"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     flashCookie,
		Value:    url.QueryEscape(kind + ":" + msg),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60,
	})
}

// getFlash reads and clears the flash cookie, returning the message and
// whether it represents an error. The cookie is cleared so the message is only
// displayed once.
func getFlash(w http.ResponseWriter, r *http.Request) (string, bool) {
	c, err := r.Cookie(flashCookie)
	if err != nil {
		return "", false
	}

	// Clear the cookie regardless of its contents.
	http.SetCookie(w, &http.Cookie{
		Name:     flashCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	value, err := url.QueryUnescape(c.Value)
	if err != nil {
		value = c.Value
	}
	kind, msg, _ := strings.Cut(value, ":")
	return msg, kind == "error"
}
