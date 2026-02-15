# Task: connectionInfo Service

## Feature Name
connectionInfo

## Description
A lightweight NixOS service that hosts a webpage displaying client connection information.

## Requirements
- Single HTTP endpoint at / that returns an HTML page
- Display client IP address (handle X-Forwarded-For for proxied requests)
- Display all request headers in a readable format
- Display request method, path, and query parameters
- Display server timestamp in UTC
- Parse and display User-Agent details (browser, OS) when possible
- Clean, minimal HTML styling with no external dependencies
- Should run on port 8080 by default, configurable via NixOS option

## Context
- This is for a NixOS flake-based system like nrvps
- Implementation should be in Go
- Deploy as a systemd service via NixOS module
- Keep it lightweight for VPS deployment
