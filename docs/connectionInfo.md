# connectionInfo Service

A lightweight diagnostic web service that displays detailed information about incoming HTTP connections. Designed for deployment on NixOS-based VPS systems.

## Overview

When you visit the connectionInfo service, it displays:

- **Your IP Address** - Shows your public IP, correctly handling requests through reverse proxies
- **Request Details** - HTTP method, path, and query parameters
- **Browser Information** - Parsed browser name, version, and operating system
- **Request Headers** - All HTTP headers in alphabetical order
- **Server Timestamp** - Current server time in UTC (ISO 8601 format)

## Quick Start

### NixOS Installation

Add the connectionInfo flake to your NixOS configuration:

```nix
{
  inputs.connectionInfo.url = "path:./src/connectionInfo";
  # Or from a remote source:
  # inputs.connectionInfo.url = "github:user/connectionInfo";

  outputs = { self, nixpkgs, connectionInfo, ... }: {
    nixosConfigurations.myhost = nixpkgs.lib.nixosSystem {
      # ...
      modules = [
        connectionInfo.nixosModules.default
        {
          services.connectionInfo = {
            enable = true;
            nginx.virtualHost = "ilios.dev";
          };
        }
      ];
    };
  };
}
```

This is all you need. Enabling the service automatically:

- Starts the connectionInfo server on port 8080 (internal only)
- Enables nginx with `recommendedProxySettings`
- Serves the page at `HOSTNAME/connectionInfo` (e.g., `ilios.dev/connectionInfo`)

### Accessing the Service

Once running, navigate to `http://your-server/connectionInfo` in your browser.

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `services.connectionInfo.enable` | boolean | `false` | Enable the connectionInfo service |
| `services.connectionInfo.port` | port | `8080` | Port for the internal server to listen on |
| `services.connectionInfo.openFirewall` | boolean | `false` | Open the firewall for the configured port (not needed when using the built-in nginx) |
| `services.connectionInfo.package` | package | (default) | The connectionInfo package to use |
| `services.connectionInfo.basePath` | string | `"/connectionInfo"` | URL path prefix where the service is hosted (empty string = serve at virtual host root) |
| `services.connectionInfo.nginx.enable` | boolean | `true` | Enable the built-in nginx reverse proxy (enabled by default) |
| `services.connectionInfo.nginx.virtualHost` | string | `"localhost"` | nginx virtual host name under which to serve the service |
| `services.connectionInfo.nginx.forceSSL` | boolean | `false` | Force SSL for the virtual host |
| `services.connectionInfo.nginx.enableACME` | boolean | `false` | Enable ACME (Let's Encrypt) certificate provisioning for the virtual host |

### Example Configurations

**Minimal setup (serves at `localhost/connectionInfo`):**

```nix
services.connectionInfo = {
  enable = true;
};
```

**Production deployment with TLS on a public domain:**

```nix
services.connectionInfo = {
  enable = true;
  nginx.virtualHost = "ilios.dev";
  nginx.forceSSL = true;
  nginx.enableACME = true;
};

# Let's Encrypt / ACME certificate provisioning
security.acme = {
  acceptTerms = true;
  defaults.email = "you@example.com";
};
```

This serves the page at `https://ilios.dev/connectionInfo` with automatic TLS.

**Custom base path:**

```nix
services.connectionInfo = {
  enable = true;
  basePath = "/diagnostics";
  nginx.virtualHost = "ilios.dev";
};
```

This serves the page at `ilios.dev/diagnostics` instead of the default `/connectionInfo`.

**Serve at the virtual host root (no base path):**

```nix
services.connectionInfo = {
  enable = true;
  basePath = "";
  nginx.virtualHost = "connectioninfo.ilios.dev";
};
```

This serves the page directly at `connectioninfo.ilios.dev/`.

**Disable the built-in nginx (manage your own reverse proxy):**

```nix
services.connectionInfo = {
  enable = true;
  nginx.enable = false;
  openFirewall = true;
};
```

## How the Built-in nginx Works

By default, enabling connectionInfo also configures nginx automatically. The module:

1. Enables `services.nginx` with `recommendedProxySettings` (which sets `X-Forwarded-For`, `X-Real-IP`, and other proxy headers)
2. Creates a virtual host entry under `services.connectionInfo.nginx.virtualHost`
3. Configures a location block at `basePath` that proxies to the internal connectionInfo server

When `basePath` is set (the default is `/connectionInfo`), nginx is configured to:
- Redirect `HOSTNAME/connectionInfo` → `HOSTNAME/connectionInfo/` (301)
- Proxy `HOSTNAME/connectionInfo/*` to the internal server with the prefix stripped

This means the connectionInfo server always receives requests at `/` regardless of the base path — nginx handles the rewriting transparently.

> **Note:** Since nginx is the public-facing listener, you typically do **not** need `openFirewall = true`. The internal server listens on `127.0.0.1` only via the configured port.

## API Reference

### GET /

Returns an HTML page displaying connection information.

**Response:**
- Status: `200 OK`
- Content-Type: `text/html; charset=utf-8`

**Response sections:**
- Your IP Address
- Request Details (method, path, query parameters)
- Your Browser (parsed browser/OS info + raw User-Agent)
- Request Headers (alphabetically sorted)
- Server Timestamp (UTC, ISO 8601)

### All Other Paths

Any path other than the root (after nginx prefix stripping) returns a 404 response.

**Response:**
- Status: `404 Not Found`
- Content-Type: `text/plain`
- Body: `404 Not Found`

## Understanding the Output

### IP Address Display

The service determines your IP address using this logic:

1. If `X-Forwarded-For` header exists, the **leftmost** IP is displayed (original client)
2. Otherwise, the direct connection IP is displayed

### Browser Detection

The service parses the User-Agent header to detect:

**Supported browsers:**
- Chrome
- Firefox
- Safari
- Edge
- Opera

**Supported operating systems:**
- Windows (10, 11)
- macOS
- Linux
- iOS
- Android
- ChromeOS

For unrecognized User-Agents, "Unknown" is displayed but the raw User-Agent string is always shown.

### Query Parameters

Query parameters are URL-decoded and displayed. For example:

```
http://server:8080/?name=John%20Doe&debug=true
```

Displays as:
- `name`: `John Doe`
- `debug`: `true`

## Troubleshooting

### Service won't start

Check the systemd service status:

```bash
systemctl status connectionInfo
journalctl -u connectionInfo -f
```

### Wrong IP address displayed

If you see `127.0.0.1` or your proxy's IP instead of the client IP:

1. Ensure `services.connectionInfo.nginx.enable` is `true` (the default) — the built-in nginx config sets `X-Forwarded-For` automatically via `recommendedProxySettings`
2. If using a custom reverse proxy, ensure it sets the `X-Forwarded-For` header

### Port already in use

Change the port in your NixOS configuration:

```nix
services.connectionInfo.port = 8081;
```

### Can't access from external network

Enable the firewall option:

```nix
services.connectionInfo.openFirewall = true;
```

Or manually open the port in your firewall configuration.

## Security Notes

- The service runs with systemd security hardening (DynamicUser, NoNewPrivileges, ProtectSystem, etc.)
- All user input is HTML-escaped to prevent XSS attacks
- No authentication is provided - the service is intended for diagnostic purposes
- TLS/HTTPS is handled by the built-in nginx reverse proxy, not by the service itself
- The `X-Forwarded-For` header is trusted unconditionally - the built-in nginx ensures this is set correctly; if using a custom proxy, deploy behind a trusted one

## Known Limitations

- **User-Agent parsing**: Limited to top 5 browsers (Chrome, Firefox, Safari, Edge, Opera)
- **Windows 11 detection**: Uses Win64 heuristic which may not be 100% accurate in all cases
- **X-Forwarded-For trust**: The header is trusted unconditionally; the built-in nginx handles this correctly, but custom proxy setups must ensure the header is trustworthy
- **Single endpoint**: Only the root path (`/`) is served; all other paths return 404 (nginx handles base path rewriting)
- **HTML only**: No JSON API endpoint is provided
- **Shared virtual host**: When using `basePath`, the nginx virtual host is configured with `lib.mkMerge`, allowing other services to add their own locations to the same virtual host
