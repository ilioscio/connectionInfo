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
            port = 8080;
            openFirewall = true;
          };
        }
      ];
    };
  };
}
```

### Accessing the Service

Once running, navigate to `http://your-server:8080/` in your browser.

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `services.connectionInfo.enable` | boolean | `false` | Enable the connectionInfo service |
| `services.connectionInfo.port` | port | `8080` | Port for the server to listen on |
| `services.connectionInfo.openFirewall` | boolean | `false` | Open the firewall for the configured port |
| `services.connectionInfo.package` | package | (default) | The connectionInfo package to use |

### Example Configurations

**Basic setup on port 8080:**

```nix
services.connectionInfo = {
  enable = true;
};
```

**Custom port with firewall:**

```nix
services.connectionInfo = {
  enable = true;
  port = 3000;
  openFirewall = true;
};
```

## Deployment with Reverse Proxy

The service is designed to run behind a reverse proxy (nginx, Caddy, etc.) that handles TLS termination. The service reads the `X-Forwarded-For` header to determine the original client IP.

### nginx Example

```nginx
server {
    listen 443 ssl;
    server_name connectioninfo.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Caddy Example

```caddy
connectioninfo.example.com {
    reverse_proxy localhost:8080
}
```

Caddy automatically sets `X-Forwarded-For` headers.

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

Returns a 404 response.

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

1. Ensure your reverse proxy sets the `X-Forwarded-For` header
2. For nginx, add: `proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;`

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
- TLS/HTTPS should be handled by a reverse proxy, not by this service
- The `X-Forwarded-For` header is trusted unconditionally - deploy behind a trusted reverse proxy

## Known Limitations

- **User-Agent parsing**: Limited to top 5 browsers (Chrome, Firefox, Safari, Edge, Opera)
- **Windows 11 detection**: Uses Win64 heuristic which may not be 100% accurate in all cases
- **X-Forwarded-For trust**: The header is trusted unconditionally, suitable only for deployment behind trusted reverse proxies
- **Single endpoint**: Only the root path (`/`) is served; all other paths return 404
- **HTML only**: No JSON API endpoint is provided
