
# gosqlpp MCP Proxy

[![Go Version](https://img.shields.io/badge/go-1.24.5-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust Go-based proxy utility for forwarding and logging Model Context Protocol (MCP) traffic between clients and the sqlpp MCP server. This tool is essential for debugging, auditing, and monitoring MCP communications in development and production environments.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [Logging](#logging)
- [Development](#development)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Overview

The Model Context Protocol (MCP) enables standardized communication between applications and context providers. This proxy sits between MCP clients and the sqlpp MCP server, providing transparent traffic forwarding while capturing all communications for analysis and debugging.

**Key Use Cases:**
- Debugging MCP client-server interactions
- Auditing database queries and responses
- Monitoring performance and traffic patterns
- Development and testing of MCP integrations

## Features

- **Dual Transport Support**: Both stdio and HTTP transport modes
- **Flexible Configuration**: Command-line flags, environment variables, and config files (YAML/JSON/TOML)
- **Configurable Executable Path**: Specify the path to mcp_sqlpp executable via flag or config
- **Comprehensive Logging**: All traffic logged to unique files per run with timestamps
- **Zero-Configuration**: Works out of the box with sensible defaults
- **Production Ready**: Robust error handling and graceful shutdown
- **Cross-Platform**: Runs on macOS, Linux, and Windows

## Prerequisites

- **Go 1.24.5 or later** (for building from source)
- **sqlpp MCP server** installed and accessible
- **Network access** to the target sqlpp server (for HTTP mode)

## Installation

### Option 1: Build from Source
```bash
# Clone the repository
git clone https://github.com/your-org/gosqlpp-mcp-proxy.git
cd gosqlpp-mcp-proxy

# Build the binary
go build -o mcp_sqlpp_proxy main.go

# Verify the build
./mcp_sqlpp_proxy --help
```

### Option 2: Install with Go
```bash
go install github.com/your-org/gosqlpp-mcp-proxy@latest
```

## Quick Start

### Setup sqlpp Path
The proxy requires access to the sqlpp MCP server executable. You can specify the path in several ways:

1. **Using command-line flag (recommended):**
   ```bash
   ./mcp_sqlpp_proxy --exe-path /path/to/your/mcp_sqlpp
   ```

2. **Using configuration file:**
   ```yaml
   # config.yaml
   exe-path: /path/to/your/mcp_sqlpp
   ```

3. **Using environment variable:**
   ```bash
   export EXE_PATH=/path/to/your/mcp_sqlpp
   ./mcp_sqlpp_proxy
   ```

4. **Find your sqlpp installation:**
   ```bash
   which mcp_sqlpp
   # or search common locations
   find /usr /opt ~ -name "mcp_sqlpp" 2>/dev/null
   ```

### 1. Stdio Mode (Default)
Perfect for command-line integrations and simple debugging:

```bash
# Using default executable path (./mcp_sqlpp)
./mcp_sqlpp_proxy --transport stdio

# Specifying custom executable path
./mcp_sqlpp_proxy --transport stdio --exe-path /usr/local/bin/mcp_sqlpp
```

### 2. HTTP Mode
Ideal for web applications and HTTP-based integrations:

```bash
# Basic HTTP proxy
./mcp_sqlpp_proxy --transport http --port 8080 --xfer-port 8891

# With custom executable path
./mcp_sqlpp_proxy --transport http --port 8080 --xfer-port 8891 --exe-path /usr/local/bin/mcp_sqlpp
```

### 3. With Configuration File
For complex setups and production deployments:

```bash
./mcp_sqlpp_proxy --config config.yaml
```

## Configuration

### Command-Line Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--transport` | `-t` | `stdio` | Transport mode: `stdio` or `http` |
| `--port` | `-p` | `8099` | Port to listen on (HTTP mode only) |
| `--xfer-port` | `-x` | `8891` | Port where sqlpp MCP server is running |
| `--exe-path` | `-e` | `./mcp_sqlpp` | Path to the mcp_sqlpp executable |
| `--config` | | | Path to configuration file |
| `--help` | `-h` | | Show help message |

### Environment Variables

All configuration options can be set via environment variables with the `MCP_PROXY_` prefix:

```bash
export MCP_PROXY_TRANSPORT=http
export MCP_PROXY_PORT=8080
export MCP_PROXY_XFER_PORT=8891
export MCP_PROXY_EXE_PATH=/usr/local/bin/mcp_sqlpp
./mcp_sqlpp_proxy
```

### Configuration File

The proxy supports multiple configuration file formats and locations:

**File Names (searched in order):**
- `mcp_sqlpp_proxy.yaml` (recommended)
- `mcp_sqlpp_proxy.json`
- `mcp_sqlpp_proxy.toml`

**Search Paths:**
- Current directory (`.`)
- `./config/`
- `/etc/mcp-proxy/`

**Or specify a custom config file:**
```bash
./mcp_sqlpp_proxy --config /path/to/custom-config.yaml
```

#### YAML Configuration (Recommended)
```yaml
# mcp_sqlpp_proxy.yaml
transport: http
port: 8080
xfer-port: 8891
exe-path: /usr/local/bin/mcp_sqlpp
```

#### JSON Configuration
```json
{
  "transport": "http",
  "port": 8080,
  "xfer-port": 8891,
  "exe-path": "/usr/local/bin/mcp_sqlpp"
}
```

#### TOML Configuration
```toml
transport = "http"
port = 8080
xfer-port = 8891
exe-path = "/usr/local/bin/mcp_sqlpp"
```

### Configuration Precedence

Configuration values are applied in the following order (highest to lowest priority):

1. **Command-line flags** (highest priority)
2. **Environment variables** (`MCP_PROXY_*`)
3. **Configuration file**
4. **Built-in defaults** (lowest priority)

### Default Configuration File

Use the included example configuration:

```bash
# Copy the example config
cp mcp_sqlpp_proxy.yaml my-config.yaml

# Edit as needed
vim my-config.yaml

# Use your custom config
./mcp_sqlpp_proxy --config my-config.yaml
```

## Usage Examples

### Development Debugging
Monitor all MCP traffic during development:

```bash
# Terminal 1: Start the proxy with custom executable path
./mcp_sqlpp_proxy --transport stdio --exe-path /usr/local/bin/mcp_sqlpp

# Terminal 2: Monitor logs in real-time
tail -f mcp_sqlpp_proxy_*.log
```

### Production HTTP Proxy
Set up a production HTTP proxy with custom ports:

```bash
./mcp_sqlpp_proxy \
  --transport http \
  --port 9000 \
  --xfer-port 8891
```

### Docker Integration
```dockerfile
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mcp_sqlpp_proxy main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mcp_sqlpp_proxy .
CMD ["./mcp_sqlpp_proxy", "--transport", "http", "--port", "8080"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mcp-proxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mcp-proxy
  template:
    metadata:
      labels:
        app: mcp-proxy
    spec:
      containers:
      - name: mcp-proxy
        image: mcp-sqlpp-proxy:latest
        ports:
        - containerPort: 8080
        env:
        - name: TRANSPORT
          value: "http"
        - name: PORT
          value: "8080"
        - name: XFER_PORT
          value: "8891"
```

## Logging

### Log File Format
Each run creates a unique log file: `mcp_sqlpp_proxy_<pid>_<timestamp>.log`

**Example filename:** `mcp_sqlpp_proxy_12345_1704067200000000000.log`

### Log Content
- **Stdio Mode**: All input/output messages with `[IN]` and `[OUT]` prefixes
- **HTTP Mode**: HTTP requests/responses with `[HTTP IN]`, `[HTTP OUT]`, and `[HTTP ERROR]` prefixes
- **Timestamps**: All entries include precise timestamps
- **Request/Response Bodies**: Full content logged for debugging

### Log Analysis
```bash
# View recent logs
ls -la mcp_sqlpp_proxy_*.log | tail -5

# Search for errors
grep "ERROR" mcp_sqlpp_proxy_*.log

# Monitor live traffic
tail -f mcp_sqlpp_proxy_$(pgrep mcp_sqlpp_proxy)_*.log
```

## Development

### Project Structure
```
gosqlpp-mcp-proxy/
├── main.go                    # Main application code
├── main_test.go              # Basic tests
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── README.md                 # This file
├── CONTRIBUTING.md           # Contribution guidelines
├── LICENSE                   # MIT license
├── .gitignore               # Git ignore rules
├── config.example.yaml      # Example configuration file
├── docs/
│   └── product.md           # Detailed product documentation
├── .github/
│   └── copilot-instructions.md
└── .vscode/                 # VS Code settings
```

### Building
```bash
# Build for current platform
go build -o mcp_sqlpp_proxy main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o mcp_sqlpp_proxy-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o mcp_sqlpp_proxy-windows-amd64.exe main.go
GOOS=darwin GOARCH=amd64 go build -o mcp_sqlpp_proxy-darwin-amd64 main.go
```

### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build test
go build -o test_binary main.go && rm test_binary
```

### Dependencies
- **Viper**: Configuration management
- **pflag**: POSIX-compliant command-line flags
- **Standard Library**: HTTP server, process management, logging

## Contributing

We welcome contributions! Please follow these guidelines:

### Getting Started
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Code Style
- Follow standard Go formatting: `go fmt ./...`
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

### Reporting Issues
Please use the GitHub issue tracker to report bugs or request features. Include:
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant log output

## Troubleshooting

### Known Issues

**HTTP Request Body Forwarding**
The current HTTP proxy implementation has a limitation where the request body is read for logging but not properly recreated for forwarding to the target server. This may cause issues with POST/PUT requests that include a body.

**Workaround**: For now, use stdio mode for requests that include a body, or consider contributing a fix for this issue.

### Common Issues

**Issue: "Failed to start mcp_sqlpp"**
```bash
# Check if sqlpp is installed and accessible
which mcp_sqlpp
# Use the --exe-path flag to specify the correct path
./mcp_sqlpp_proxy --exe-path /path/to/mcp_sqlpp
```

**Issue: "Port already in use"**
```bash
# Check what's using the port
lsof -i :8099
# Use a different port
./mcp_sqlpp_proxy --port 8100
```

**Issue: "Connection refused" (HTTP mode)**
```bash
# Verify sqlpp server is running on the expected port
curl http://localhost:8891/health
# Check firewall settings
```

**Issue: "Failed to open log file"**
```bash
# Check write permissions in current directory
ls -la .
# Run from a directory where you have write access
```

**Issue: HTTP request body not being forwarded correctly**
The current implementation reads the request body for logging but doesn't recreate it for forwarding. This is a known limitation that may cause issues with some requests.

### Debug Mode
Enable verbose logging by modifying the log level in the source code or set environment variables:

```bash
export DEBUG=true
./mcp_sqlpp_proxy --transport stdio
```

### Performance Tuning
For high-traffic scenarios:
- Increase system file descriptor limits
- Monitor memory usage with `top` or `htop`
- Consider running multiple instances behind a load balancer
- Be aware that all traffic is logged, which may impact performance

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Additional Resources

- [Model Context Protocol Specification](https://spec.modelcontextprotocol.io/)
- [sqlpp MCP Server Documentation](https://github.com/sqlpp/mcp-server)
- [Product Documentation](docs/product.md)

## Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/your-org/gosqlpp-mcp-proxy/issues)
- **Discussions**: [Community discussions and Q&A](https://github.com/your-org/gosqlpp-mcp-proxy/discussions)

---

**Made with ❤️ for the MCP community**
