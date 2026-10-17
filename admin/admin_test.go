package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Doolox",
		"Generate website",
		"app-wrapper",
		"sidebar-brand",
		"static website builder",
		"/assets/favicon.png",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index body missing %q", want)
		}
	}
}

func TestLogo(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logo.svg", nil)
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "image/svg+xml" {
		t.Errorf("content type = %q, want image/svg+xml", got)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Error("logo body missing <svg")
	}
}

func TestUnknownPathIsNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/no-such-page", nil)
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestIndexShowsFlashMessage(t *testing.T) {
	rec := httptest.NewRecorder()
	setFlash(rec, "hello", false)
	c := flashCookieOf(t, rec)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(c)
	out := httptest.NewRecorder()
	handler().ServeHTTP(out, req)

	if out.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", out.Code, http.StatusOK)
	}

	body := out.Body.String()
	if !strings.Contains(body, "hello") {
		t.Error("index body missing flash message")
	}
	if !strings.Contains(body, "alert-success") {
		t.Error("index body missing success alert")
	}
}

func TestIndexShowsFlashError(t *testing.T) {
	rec := httptest.NewRecorder()
	setFlash(rec, "boom", true)
	c := flashCookieOf(t, rec)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(c)
	out := httptest.NewRecorder()
	handler().ServeHTTP(out, req)

	body := out.Body.String()
	if !strings.Contains(body, "boom") {
		t.Error("index body missing flash message")
	}
	if !strings.Contains(body, "alert-danger") {
		t.Error("index body missing error alert")
	}
}

func TestFlashShownOnceAndCleared(t *testing.T) {
	rec := httptest.NewRecorder()
	setFlash(rec, "once", false)
	c := flashCookieOf(t, rec)

	// The first request displays the message and clears the cookie.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(c)
	out := httptest.NewRecorder()
	handler().ServeHTTP(out, req)
	if !strings.Contains(out.Body.String(), "once") {
		t.Error("flash message not shown")
	}
	if got := cookieNamed(out, flashCookie); got == nil || got.MaxAge != -1 {
		t.Error("flash cookie not cleared after display")
	}

	// A second request without the cookie shows nothing.
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	out2 := httptest.NewRecorder()
	handler().ServeHTTP(out2, req2)
	if strings.Contains(out2.Body.String(), "once") {
		t.Error("flash message shown twice")
	}
}

func flashCookieOf(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	if c := cookieNamed(rec, flashCookie); c != nil {
		return c
	}
	t.Fatal("no flash cookie set")
	return nil
}

func cookieNamed(rec *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"#define MyAppVersion \"1.0.1\"\n", "1.0.1"},
		{"#define MyAppVersion \"2.3.4\"", "2.3.4"},
		{"ignored\n#define MyAppVersion \"0.9.0\"\ntail", "0.9.0"},
		{"no version here", ""},
	}

	for _, tt := range tests {
		if got := parseVersion([]byte(tt.in)); got != tt.want {
			t.Errorf("parseVersion(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
