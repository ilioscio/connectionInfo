# Handoff: Architecture → Engineering

## Task Reference
- Feature: connectionInfo
- Requirements: requirements/connectionInfo.md
- Design: design/connectionInfo.md

## Completed Work

### Package Structure Defined
- Single Go module with `internal/` packages for server, handler, parser, and render
- Entry point in `main.go`
- Flake.nix with package, NixOS module, and overlay

### Technical Decisions Made
1. **Go standard library only** - No external dependencies for User-Agent parsing or HTTP; use regexp and html/template
2. **X-Forwarded-For trust** - Trust unconditionally, take leftmost IP (this is a diagnostic tool, not security-critical)
3. **HTML rendering** - Use embedded html/template with automatic escaping
4. **User-Agent parsing** - Regex patterns for top 5 browsers and major OS; fallback to "Unknown" gracefully
5. **Error responses** - Minimal responses (404 plain text, no stack traces)

### Data Models Specified
- `ConnectionInfo` struct for passing data to renderer
- `HeaderPair` for sorted header display
- `UserAgentInfo` for parsed browser/OS data

### NixOS Module Options
- `enable` - boolean to enable service
- `port` - defaults to 8080
- `openFirewall` - optional firewall rule

### Security Hardening
- Full systemd hardening config specified (DynamicUser, ProtectSystem=strict, etc.)
- All user input HTML-escaped before rendering

## Handoff To
Engineering

## Implementation Order (Recommended)
1. Set up Go module and basic flake.nix structure
2. Implement `internal/parser/ip.go` with X-Forwarded-For parsing
3. Implement `internal/parser/useragent.go` with regex-based parsing
4. Implement `internal/render/html.go` with template
5. Implement `internal/handler/handler.go` to tie it together
6. Implement `internal/server/server.go` and `main.go`
7. Complete flake.nix with NixOS module and systemd config
8. Write unit tests for parsers and renderer

## Questions/Blockers

**None** - All technical concerns from Product handoff have been addressed:
- X-Forwarded-For: Trust unconditionally (documented rationale)
- User-Agent parsing depth: Regex for top browsers, fallback gracefully
- Error handling: Minimal responses, no information leakage
