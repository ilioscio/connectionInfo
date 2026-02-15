package render

import (
	"html/template"
	"io"
	"time"

	"connectionInfo/internal/parser"
)

// ConnectionInfo holds all data to be rendered in the HTML page.
type ConnectionInfo struct {
	ClientIP      string
	RawRemoteAddr string
	Method        string
	Path          string
	QueryParams   map[string][]string
	Headers       []HeaderPair
	UserAgent     parser.UserAgentInfo
	Timestamp     time.Time
}

// HeaderPair represents a single HTTP header key-value pair.
type HeaderPair struct {
	Name  string
	Value string
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Connection Info</title>
    <style>
        * {
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
            color: #333;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        h2 {
            color: #34495e;
            margin-top: 30px;
            font-size: 1.2em;
        }
        section {
            background: white;
            padding: 20px;
            margin: 20px 0;
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .ip-address {
            font-size: 2em;
            font-weight: bold;
            color: #3498db;
            font-family: monospace;
            margin: 10px 0;
        }
        dl {
            display: grid;
            grid-template-columns: auto 1fr;
            gap: 8px 16px;
            margin: 0;
        }
        dt {
            font-weight: 600;
            color: #555;
        }
        dd {
            margin: 0;
            font-family: monospace;
            word-break: break-all;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 0.9em;
        }
        th, td {
            text-align: left;
            padding: 8px 12px;
            border-bottom: 1px solid #eee;
        }
        th {
            background: #f8f9fa;
            font-weight: 600;
            color: #555;
        }
        td:first-child {
            font-weight: 500;
            white-space: nowrap;
        }
        td:last-child {
            font-family: monospace;
            word-break: break-all;
        }
        .timestamp {
            font-family: monospace;
            color: #666;
        }
        .raw-ua {
            font-size: 0.85em;
            color: #666;
            word-break: break-all;
        }
        @media (max-width: 600px) {
            body {
                padding: 10px;
            }
            dl {
                grid-template-columns: 1fr;
            }
            dt {
                margin-top: 10px;
            }
            .ip-address {
                font-size: 1.5em;
            }
        }
    </style>
</head>
<body>
    <h1>Connection Information</h1>

    <section id="ip">
        <h2>Your IP Address</h2>
        <p class="ip-address">{{.ClientIP}}</p>
    </section>

    <section id="request">
        <h2>Request Details</h2>
        <dl>
            <dt>Method</dt>
            <dd>{{.Method}}</dd>
            <dt>Path</dt>
            <dd>{{.Path}}</dd>
            <dt>Query Parameters</dt>
            <dd>{{if .QueryParams}}{{range $key, $values := .QueryParams}}{{$key}}={{range $i, $v := $values}}{{if $i}}, {{end}}{{$v}}{{end}} {{end}}{{else}}(none){{end}}</dd>
        </dl>
    </section>

    <section id="useragent">
        <h2>Your Browser</h2>
        <dl>
            <dt>Browser</dt>
            <dd>{{.UserAgent.BrowserName}}{{if .UserAgent.BrowserVersion}} {{.UserAgent.BrowserVersion}}{{end}}</dd>
            <dt>Operating System</dt>
            <dd>{{.UserAgent.OSName}}</dd>
            <dt>Raw User-Agent</dt>
            <dd class="raw-ua">{{if .UserAgent.Raw}}{{.UserAgent.Raw}}{{else}}(not provided){{end}}</dd>
        </dl>
    </section>

    <section id="headers">
        <h2>Request Headers</h2>
        <table>
            <thead>
                <tr>
                    <th>Header</th>
                    <th>Value</th>
                </tr>
            </thead>
            <tbody>
                {{range .Headers}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.Value}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </section>

    <section id="timestamp">
        <h2>Server Timestamp</h2>
        <p class="timestamp">{{.Timestamp.Format "2006-01-02T15:04:05Z07:00"}}</p>
    </section>
</body>
</html>`

var tmpl = template.Must(template.New("connectionInfo").Parse(htmlTemplate))

// Render writes the HTML page to the provided writer.
func Render(w io.Writer, info ConnectionInfo) error {
	return tmpl.Execute(w, info)
}
