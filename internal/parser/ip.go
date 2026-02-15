package parser

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP extracts the client IP address from the request.
// It first checks the X-Forwarded-For header (taking the leftmost IP),
// then falls back to the RemoteAddr from the request.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// The leftmost IP is the original client
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if clientIP != "" {
				return clientIP
			}
		}
	}

	// Check X-Real-IP header (commonly used by nginx)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	return extractIP(r.RemoteAddr)
}

// extractIP extracts the IP address from a host:port string.
// If the input is already just an IP, it returns it unchanged.
func extractIP(remoteAddr string) string {
	// RemoteAddr is typically "IP:port" or "[IPv6]:port"
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// Might be just an IP without port
		return remoteAddr
	}
	return host
}
