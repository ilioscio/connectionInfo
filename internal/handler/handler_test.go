package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_RootPath(t *testing.T) {
	h := New()

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/html; charset=utf-8")
	}

	body := rr.Body.String()

	// Check for expected content sections
	expectedContents := []string{
		"Connection Information",
		"Your IP Address",
		"192.168.1.100",
		"Request Details",
		"GET",
		"Your Browser",
		"Chrome",
		"Windows",
		"Request Headers",
		"Server Timestamp",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(body, expected) {
			t.Errorf("response body does not contain %q", expected)
		}
	}
}

func TestHandler_NonRootPath_Returns404(t *testing.T) {
	h := New()

	paths := []string{
		"/api",
		"/health",
		"/favicon.ico",
		"/something/else",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("handler returned wrong status code for %s: got %v want %v", path, status, http.StatusNotFound)
			}

			body := rr.Body.String()
			if !strings.Contains(body, "404 Not Found") {
				t.Errorf("response body should contain '404 Not Found', got %q", body)
			}
		})
	}
}

func TestHandler_XForwardedFor(t *testing.T) {
	h := New()

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.50, 70.41.3.18")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "203.0.113.50") {
		t.Errorf("response should contain client IP from X-Forwarded-For, got %q", body)
	}
}

func TestHandler_QueryParams(t *testing.T) {
	h := New()

	req := httptest.NewRequest("GET", "/?foo=bar&baz=qux", nil)
	req.RemoteAddr = "192.168.1.100:12345"

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	// Check that query params are displayed
	if !strings.Contains(body, "foo") || !strings.Contains(body, "bar") {
		t.Errorf("response should contain query parameters, got %q", body)
	}
}

func TestHandler_NoUserAgent(t *testing.T) {
	h := New()

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	// Don't set User-Agent

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Unknown") {
		t.Errorf("response should show 'Unknown' for missing User-Agent, got %q", body)
	}
}

func TestHandler_HTMLEscaping(t *testing.T) {
	h := New()

	req := httptest.NewRequest("GET", "/?xss=<script>alert('xss')</script>", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	req.Header.Set("X-Custom-Header", "<script>malicious</script>")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()

	// Should NOT contain unescaped script tags
	if strings.Contains(body, "<script>") {
		t.Errorf("response should not contain unescaped script tags, got %q", body)
	}

	// Should contain escaped versions
	if !strings.Contains(body, "&lt;script&gt;") {
		t.Errorf("response should contain HTML-escaped script tags, got %q", body)
	}
}

func TestHandler_MethodsDisplayed(t *testing.T) {
	h := New()

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			req.RemoteAddr = "192.168.1.100:12345"

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			body := rr.Body.String()
			if !strings.Contains(body, method) {
				t.Errorf("response should contain method %s, got %q", method, body)
			}
		})
	}
}

func TestHandler_HeadersSorted(t *testing.T) {
	h := New()

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	req.Header.Set("Zebra-Header", "z-value")
	req.Header.Set("Alpha-Header", "a-value")
	req.Header.Set("Middle-Header", "m-value")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	body := rr.Body.String()

	// Find positions of headers in the output
	alphaPos := strings.Index(body, "Alpha-Header")
	middlePos := strings.Index(body, "Middle-Header")
	zebraPos := strings.Index(body, "Zebra-Header")

	if alphaPos == -1 || middlePos == -1 || zebraPos == -1 {
		t.Errorf("all headers should be present in output")
		return
	}

	if !(alphaPos < middlePos && middlePos < zebraPos) {
		t.Errorf("headers should be sorted alphabetically")
	}
}
