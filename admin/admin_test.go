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
	req := httptest.NewRequest(http.MethodGet, "/?status=success&msg=hello", nil)
	rec := httptest.NewRecorder()

	handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "hello") {
		t.Error("index body missing flash message")
	}
	if !strings.Contains(body, "alert-success") {
		t.Error("index body missing success alert")
	}
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
