package cms

import "testing"

func TestParseFrontMatter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]string
	}{
		{
			name: "standard block",
			in:   "<!--\ntitle: About Doolox\nsubtitle: About\n-->",
			want: map[string]string{
				"title":    "About Doolox",
				"subtitle": "About",
			},
		},
		{
			name: "whitespace and case",
			in:   "<!--\n  Title:  Hello  \nTYPE: page\n-->",
			want: map[string]string{
				"title": "Hello",
				"type":  "page",
			},
		},
		{
			name: "ignores malformed lines",
			in:   "<!--\ntitle: X\nno colon here\n: empty key\n-->",
			want: map[string]string{"title": "X"},
		},
		{
			name: "no front matter",
			in:   "{{define \"about\"}}hello{{end}}",
			want: nil,
		},
		{
			name: "leading whitespace then comment",
			in:   "\n  <!--\ntitle: X\n-->\n{{define \"x\"}}{{end}}",
			want: map[string]string{"title": "X"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFrontMatter([]byte(tt.in))
			if len(got) != len(tt.want) {
				t.Fatalf("parseFrontMatter() = %v, want %v", got, tt.want)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("parseFrontMatter()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	data := []byte("<!--\ntitle: About Doolox\nsubtitle: About\ndescription: About Doolox\ntype: page\n-->\n{{define \"about\"}}...{{end}}")

	p := newPage("about", data)
	if p.URL != "about" {
		t.Errorf("URL = %q, want %q", p.URL, "about")
	}
	if p.Title != "About Doolox" {
		t.Errorf("Title = %q, want %q", p.Title, "About Doolox")
	}
	if p.Subtitle != "About" {
		t.Errorf("Subtitle = %q, want %q", p.Subtitle, "About")
	}
	if p.Description != "About Doolox" {
		t.Errorf("Description = %q, want %q", p.Description, "About Doolox")
	}
	if p.Type != "page" {
		t.Errorf("Type = %q, want %q", p.Type, "page")
	}
}
