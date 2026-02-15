# Handoff: Engineering → QA

## Task Reference
- Feature: connectionInfo
- Requirements: requirements/connectionInfo.md
- Design: design/connectionInfo.md

## Completed Work

### Files Created

**Source Code:**
- `src/connectionInfo/main.go` - Application entry point, reads PORT env var
- `src/connectionInfo/internal/server/server.go` - HTTP server setup
- `src/connectionInfo/internal/handler/handler.go` - Request handler with routing logic
- `src/connectionInfo/internal/parser/ip.go` - IP address extraction (X-Forwarded-For, X-Real-IP, RemoteAddr)
- `src/connectionInfo/internal/parser/useragent.go` - User-Agent parsing for browsers and OS
- `src/connectionInfo/internal/render/html.go` - HTML template with embedded CSS

**Test Files:**
- `src/connectionInfo/internal/parser/ip_test.go` - IP parser unit tests (9 test cases)
- `src/connectionInfo/internal/parser/useragent_test.go` - User-Agent parser tests (11 test cases)
- `src/connectionInfo/internal/render/html_test.go` - HTML rendering tests (4 test cases)
- `src/connectionInfo/internal/handler/handler_test.go` - Integration tests (8 test suites)

**Nix/Build Files:**
- `src/connectionInfo/flake.nix` - Nix flake with package, NixOS module, and overlay
- `src/connectionInfo/flake.lock` - Locked flake inputs
- `src/connectionInfo/go.mod` - Go module definition
- `src/connectionInfo/.gitignore` - Ignores Nix build artifacts

### Implementation Summary

The connectionInfo service is implemented as designed:

1. **IP Detection**: Parses X-Forwarded-For (leftmost IP), falls back to X-Real-IP, then RemoteAddr
2. **User-Agent Parsing**: Regex-based detection for Chrome, Firefox, Safari, Edge, Opera on Windows, macOS, iOS, Android, Linux, ChromeOS
3. **HTML Rendering**: Clean, responsive HTML with embedded CSS, automatic HTML escaping via Go's html/template
4. **Request Handler**: Returns 200 OK with HTML at `/`, 404 for all other paths
5. **Server**: Binds to port from PORT env var (default 8080)
6. **NixOS Module**: Configurable enable, port, openFirewall options with full systemd hardening

## Test Coverage

### Test Files
- `src/connectionInfo/internal/parser/ip_test.go` - Tests for IP extraction
- `src/connectionInfo/internal/parser/useragent_test.go` - Tests for User-Agent parsing
- `src/connectionInfo/internal/render/html_test.go` - Tests for HTML rendering/escaping
- `src/connectionInfo/internal/handler/handler_test.go` - Integration tests for the handler

### How to Run Tests

```bash
cd src/connectionInfo
nix-shell -p go --run "go test ./... -v"
```

Or with Nix:

```bash
cd src/connectionInfo
nix develop -c go test ./... -v
```

### Test Results

All 32 tests pass:
- `internal/handler`: 8 test suites (root path, 404s, XFF, query params, no UA, HTML escaping, methods, header sorting)
- `internal/parser`: 13 tests (9 IP tests, 11 UA tests)
- `internal/render`: 4 tests (basic render, HTML escaping, empty query params, empty UA)

### How to Build

```bash
cd src/connectionInfo
nix build
./result/bin/connectionInfo
```

## Handoff To
QA

## Questions/Blockers

None. Implementation follows the design document exactly.

## Known Limitations

1. **User-Agent parsing**: Only covers top 5 browsers and major operating systems. Unusual or bot user agents will show "Unknown" but the raw UA string is always displayed.

2. **X-Forwarded-For trust**: As specified in the design, the header is trusted unconditionally. If deployed directly without a proxy, clients could spoof their IP. This is documented as acceptable for a diagnostic tool.

3. **Windows 11 detection**: Windows 11 UA strings are nearly identical to Windows 10. The parser uses the Win64 hint to distinguish them, which may not be 100% accurate.

## Manual Testing Checklist for QA

- [ ] Direct connection shows RemoteAddr IP
- [ ] Behind nginx/reverse proxy, shows correct X-Forwarded-For IP
- [ ] All headers displayed alphabetically
- [ ] Chrome User-Agent correctly parsed
- [ ] Firefox User-Agent correctly parsed
- [ ] Safari User-Agent correctly parsed
- [ ] Unknown User-Agent shows "Unknown" but displays raw value
- [ ] Query parameters properly decoded and displayed
- [ ] Page renders correctly on mobile viewport
- [ ] Timestamp is in UTC and ISO 8601 format
- [ ] Non-root paths return 404
- [ ] HTML injection attempts are escaped
- [ ] Service starts on configured port
- [ ] Service restarts on failure
