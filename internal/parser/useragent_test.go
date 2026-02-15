package parser

import "testing"

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name            string
		ua              string
		wantBrowser     string
		wantVersion     string
		wantOS          string
		wantParsed      bool
	}{
		{
			name:        "Chrome on Windows",
			ua:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			wantBrowser: "Chrome",
			wantVersion: "120.0",
			wantOS:      "Windows 11",
			wantParsed:  true,
		},
		{
			name:        "Firefox on Linux",
			ua:          "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
			wantBrowser: "Firefox",
			wantVersion: "121.0",
			wantOS:      "Linux",
			wantParsed:  true,
		},
		{
			name:        "Safari on macOS",
			ua:          "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
			wantBrowser: "Safari",
			wantVersion: "17.2",
			wantOS:      "macOS",
			wantParsed:  true,
		},
		{
			name:        "Edge on Windows",
			ua:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			wantBrowser: "Edge",
			wantVersion: "120.0",
			wantOS:      "Windows 11",
			wantParsed:  true,
		},
		{
			name:        "Opera on Windows",
			ua:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
			wantBrowser: "Opera",
			wantVersion: "106.0",
			wantOS:      "Windows 11",
			wantParsed:  true,
		},
		{
			name:        "Chrome on Android",
			ua:          "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.43 Mobile Safari/537.36",
			wantBrowser: "Chrome",
			wantVersion: "120.0",
			wantOS:      "Android",
			wantParsed:  true,
		},
		{
			name:        "Safari on iOS",
			ua:          "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
			wantBrowser: "Safari",
			wantVersion: "17.2",
			wantOS:      "iOS",
			wantParsed:  true,
		},
		{
			name:        "Chrome on Windows 10",
			ua:          "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			wantBrowser: "Chrome",
			wantVersion: "120.0",
			wantOS:      "Windows 10",
			wantParsed:  true,
		},
		{
			name:        "Unknown browser",
			ua:          "CustomBot/1.0",
			wantBrowser: "Unknown",
			wantVersion: "",
			wantOS:      "Unknown",
			wantParsed:  false,
		},
		{
			name:        "Empty user agent",
			ua:          "",
			wantBrowser: "Unknown",
			wantVersion: "",
			wantOS:      "Unknown",
			wantParsed:  false,
		},
		{
			name:        "curl",
			ua:          "curl/8.4.0",
			wantBrowser: "Unknown",
			wantVersion: "",
			wantOS:      "Unknown",
			wantParsed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseUserAgent(tt.ua)

			if result.Raw != tt.ua {
				t.Errorf("Raw = %q, want %q", result.Raw, tt.ua)
			}
			if result.BrowserName != tt.wantBrowser {
				t.Errorf("BrowserName = %q, want %q", result.BrowserName, tt.wantBrowser)
			}
			if result.BrowserVersion != tt.wantVersion {
				t.Errorf("BrowserVersion = %q, want %q", result.BrowserVersion, tt.wantVersion)
			}
			if result.OSName != tt.wantOS {
				t.Errorf("OSName = %q, want %q", result.OSName, tt.wantOS)
			}
			if result.Parsed != tt.wantParsed {
				t.Errorf("Parsed = %v, want %v", result.Parsed, tt.wantParsed)
			}
		})
	}
}
