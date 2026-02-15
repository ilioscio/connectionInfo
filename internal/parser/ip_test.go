package parser

import (
	"net/http"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		expected   string
	}{
		{
			name:       "direct connection IPv4",
			remoteAddr: "192.168.1.100:12345",
			headers:    nil,
			expected:   "192.168.1.100",
		},
		{
			name:       "direct connection IPv6",
			remoteAddr: "[::1]:12345",
			headers:    nil,
			expected:   "::1",
		},
		{
			name:       "X-Forwarded-For single IP",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.50"},
			expected:   "203.0.113.50",
		},
		{
			name:       "X-Forwarded-For multiple IPs",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.50, 70.41.3.18, 150.172.238.178"},
			expected:   "203.0.113.50",
		},
		{
			name:       "X-Forwarded-For with spaces",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "  203.0.113.50  "},
			expected:   "203.0.113.50",
		},
		{
			name:       "X-Real-IP header",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Real-IP": "198.51.100.178"},
			expected:   "198.51.100.178",
		},
		{
			name:       "X-Forwarded-For takes precedence over X-Real-IP",
			remoteAddr: "10.0.0.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.50",
				"X-Real-IP":       "198.51.100.178",
			},
			expected: "203.0.113.50",
		},
		{
			name:       "empty X-Forwarded-For falls back to X-Real-IP",
			remoteAddr: "10.0.0.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "",
				"X-Real-IP":       "198.51.100.178",
			},
			expected: "198.51.100.178",
		},
		{
			name:       "RemoteAddr without port",
			remoteAddr: "192.168.1.100",
			headers:    nil,
			expected:   "192.168.1.100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			result := GetClientIP(req)
			if result != tt.expected {
				t.Errorf("GetClientIP() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.1:8080", "192.168.1.1"},
		{"[::1]:8080", "::1"},
		{"192.168.1.1", "192.168.1.1"},
		{"::1", "::1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractIP(tt.input)
			if result != tt.expected {
				t.Errorf("extractIP(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
