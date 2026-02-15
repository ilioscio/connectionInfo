package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"connectionInfo/internal/parser"
)

func TestRender(t *testing.T) {
	info := ConnectionInfo{
		ClientIP:      "192.168.1.100",
		RawRemoteAddr: "192.168.1.100:12345",
		Method:        "GET",
		Path:          "/",
		QueryParams:   map[string][]string{"test": {"value"}},
		Headers: []HeaderPair{
			{Name: "Accept", Value: "text/html"},
			{Name: "Host", Value: "localhost"},
		},
		UserAgent: parser.UserAgentInfo{
			Raw:            "Mozilla/5.0 Chrome/120.0",
			BrowserName:    "Chrome",
			BrowserVersion: "120.0",
			OSName:         "Windows 10",
			Parsed:         true,
		},
		Timestamp: time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC),
	}

	var buf bytes.Buffer
	err := Render(&buf, info)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	body := buf.String()

	expectedContents := []string{
		"<!DOCTYPE html>",
		"Connection Information",
		"192.168.1.100",
		"GET",
		"Chrome",
		"120.0",
		"Windows 10",
		"Accept",
		"text/html",
		"2024-01-15T12:30:45Z",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(body, expected) {
			t.Errorf("rendered output does not contain %q", expected)
		}
	}
}

func TestRender_HTMLEscaping(t *testing.T) {
	info := ConnectionInfo{
		ClientIP:      "<script>alert('xss')</script>",
		RawRemoteAddr: "192.168.1.100:12345",
		Method:        "GET",
		Path:          "/",
		QueryParams:   map[string][]string{"<key>": {"<value>"}},
		Headers: []HeaderPair{
			{Name: "<Header>", Value: "<script>bad</script>"},
		},
		UserAgent: parser.UserAgentInfo{
			Raw:         "<script>ua</script>",
			BrowserName: "Unknown",
			OSName:      "Unknown",
			Parsed:      false,
		},
		Timestamp: time.Now().UTC(),
	}

	var buf bytes.Buffer
	err := Render(&buf, info)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	body := buf.String()

	// Should NOT contain unescaped script tags
	if strings.Contains(body, "<script>") {
		t.Errorf("rendered output should not contain unescaped script tags")
	}

	// Should contain HTML-escaped versions
	if !strings.Contains(body, "&lt;script&gt;") {
		t.Errorf("rendered output should contain HTML-escaped script tags")
	}
}

func TestRender_EmptyQueryParams(t *testing.T) {
	info := ConnectionInfo{
		ClientIP:      "192.168.1.100",
		RawRemoteAddr: "192.168.1.100:12345",
		Method:        "GET",
		Path:          "/",
		QueryParams:   map[string][]string{},
		Headers:       []HeaderPair{},
		UserAgent: parser.UserAgentInfo{
			Raw:         "",
			BrowserName: "Unknown",
			OSName:      "Unknown",
			Parsed:      false,
		},
		Timestamp: time.Now().UTC(),
	}

	var buf bytes.Buffer
	err := Render(&buf, info)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	body := buf.String()
	if !strings.Contains(body, "(none)") {
		t.Errorf("rendered output should show '(none)' for empty query params")
	}
}

func TestRender_EmptyUserAgent(t *testing.T) {
	info := ConnectionInfo{
		ClientIP:      "192.168.1.100",
		RawRemoteAddr: "192.168.1.100:12345",
		Method:        "GET",
		Path:          "/",
		QueryParams:   map[string][]string{},
		Headers:       []HeaderPair{},
		UserAgent: parser.UserAgentInfo{
			Raw:         "",
			BrowserName: "Unknown",
			OSName:      "Unknown",
			Parsed:      false,
		},
		Timestamp: time.Now().UTC(),
	}

	var buf bytes.Buffer
	err := Render(&buf, info)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	body := buf.String()
	if !strings.Contains(body, "(not provided)") {
		t.Errorf("rendered output should show '(not provided)' for empty user agent")
	}
}
