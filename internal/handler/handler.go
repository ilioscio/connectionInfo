package handler

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"connectionInfo/internal/parser"
	"connectionInfo/internal/render"
)

// Handler handles HTTP requests for the connectionInfo service.
type Handler struct{}

// New creates a new Handler.
func New() *Handler {
	return &Handler{}
}

// ServeHTTP handles all incoming HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only serve the root path
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// Build connection info
	info := render.ConnectionInfo{
		ClientIP:      parser.GetClientIP(r),
		RawRemoteAddr: r.RemoteAddr,
		Method:        r.Method,
		Path:          r.URL.Path,
		QueryParams:   r.URL.Query(),
		Headers:       extractHeaders(r),
		UserAgent:     parser.ParseUserAgent(r.Header.Get("User-Agent")),
		Timestamp:     time.Now().UTC(),
	}

	// Set content type and render
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.Render(w, info); err != nil {
		// If rendering fails, return a generic error (no details exposed)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// extractHeaders extracts all headers from the request and returns them sorted alphabetically.
func extractHeaders(r *http.Request) []render.HeaderPair {
	var headers []render.HeaderPair

	for name, values := range r.Header {
		// Join multiple values for the same header with ", "
		value := strings.Join(values, ", ")
		headers = append(headers, render.HeaderPair{
			Name:  name,
			Value: value,
		})
	}

	// Sort alphabetically by header name
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].Name < headers[j].Name
	})

	return headers
}
