package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractText(t *testing.T) {
	cases := []struct {
		name string
		a    string
		text string
	}{
		{
			name: "valid",
			a:    `<a href="/login">Login</a>`,
			text: "Login",
		},
		{
			name: "valid: nested",
			a:    `<a href="/login">Login <strong>as admin</a></a>`,
			text: "Login as admin",
		},
		{
			name: "valid: comments",
			a:    `<a href="/login">Login <!-- as admin --></a>`,
			text: "Login",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := parse(t, c.a)
			text := extractText(a)
			if text != c.text {
				t.Errorf("extractText(%s) == %s, expected %s", c.a, text, c.text)
			}
		})
	}
}

func TestExtractHref(t *testing.T) {
	cases := []struct {
		name string
		a    string
		href string
	}{
		{
			name: "valid",
			a:    `<a href="/login">Login</a>`,
			href: "/login",
		},
		{
			name: "missing href",
			a:    `<a>Login</a>`,
			href: "",
		},
		{
			name: "other attrs",
			a:    `<a class="link">Login</a>`,
			href: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := parse(t, c.a)
			href := extractHref(a)
			if href != c.href {
				t.Errorf("extractHref(%s) == %s, expected %s", c.a, href, c.href)
			}
		})
	}
}

func parse(t *testing.T, a string) *html.Node {
	n, err := html.Parse(strings.NewReader(a))
	if err != nil {
		t.Errorf("html.Parse failed: %v", err)
		return nil
	}

	// to skip by default added tages from html.Parse
	// like <html><head></head><body></body></html>
	return n.FirstChild.FirstChild.NextSibling.FirstChild
}
