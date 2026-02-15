package parser

import (
	"regexp"
	"strings"
)

// UserAgentInfo contains parsed information from a User-Agent string.
type UserAgentInfo struct {
	Raw            string // Original User-Agent header
	BrowserName    string // e.g., "Chrome", "Firefox", "Safari"
	BrowserVersion string // e.g., "120.0"
	OSName         string // e.g., "Windows 10", "macOS", "Linux"
	Parsed         bool   // Whether parsing succeeded
}

// Browser patterns - order matters, check more specific patterns first
var browserPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"Edge", regexp.MustCompile(`Edg(?:e|A|iOS)?/(\d+(?:\.\d+)?)`)},
	{"Opera", regexp.MustCompile(`(?:OPR|Opera)[/ ](\d+(?:\.\d+)?)`)},
	{"Chrome", regexp.MustCompile(`Chrome/(\d+(?:\.\d+)?)`)},
	{"Firefox", regexp.MustCompile(`Firefox/(\d+(?:\.\d+)?)`)},
	{"Safari", regexp.MustCompile(`Version/(\d+(?:\.\d+)?).*Safari`)},
}

// OS patterns - order matters, check more specific patterns first
var osPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"Windows 11", regexp.MustCompile(`Windows NT 10\.0.*Win64`)},
	{"Windows 10", regexp.MustCompile(`Windows NT 10\.0`)},
	{"Windows 8.1", regexp.MustCompile(`Windows NT 6\.3`)},
	{"Windows 8", regexp.MustCompile(`Windows NT 6\.2`)},
	{"Windows 7", regexp.MustCompile(`Windows NT 6\.1`)},
	{"Windows", regexp.MustCompile(`Windows`)},
	{"iOS", regexp.MustCompile(`iPhone|iPad|iPod`)},         // Check iOS before macOS (iOS UA contains "like Mac OS X")
	{"macOS", regexp.MustCompile(`Mac OS X|Macintosh`)},
	{"Android", regexp.MustCompile(`Android`)},
	{"ChromeOS", regexp.MustCompile(`CrOS`)},                // Check ChromeOS before Linux (ChromeOS contains "Linux")
	{"Linux", regexp.MustCompile(`Linux`)},
}

// ParseUserAgent parses a User-Agent string and extracts browser and OS information.
func ParseUserAgent(ua string) UserAgentInfo {
	info := UserAgentInfo{
		Raw:            ua,
		BrowserName:    "Unknown",
		BrowserVersion: "",
		OSName:         "Unknown",
		Parsed:         false,
	}

	if ua == "" {
		return info
	}

	// Parse browser
	for _, bp := range browserPatterns {
		if matches := bp.pattern.FindStringSubmatch(ua); matches != nil {
			info.BrowserName = bp.name
			if len(matches) > 1 {
				info.BrowserVersion = matches[1]
			}
			info.Parsed = true
			break
		}
	}

	// Parse OS
	for _, op := range osPatterns {
		if op.pattern.MatchString(ua) {
			info.OSName = op.name
			info.Parsed = true
			break
		}
	}

	// Special case: if we detect Safari but also Chrome, it's actually Chrome
	// (Chrome includes Safari in its UA string)
	if info.BrowserName == "Safari" && strings.Contains(ua, "Chrome") {
		// Re-check for Chrome
		if matches := browserPatterns[2].pattern.FindStringSubmatch(ua); matches != nil {
			info.BrowserName = "Chrome"
			if len(matches) > 1 {
				info.BrowserVersion = matches[1]
			}
		}
	}

	return info
}
