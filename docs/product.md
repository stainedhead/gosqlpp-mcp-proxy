# gosqlpp_mcp_proxy Product Document

## Overview

gosqlpp_mcp_proxy is a Go-based proxy utility designed to forward and log Model Context Protocol (MCP) traffic between clients and the sqlpp MCP server. It supports both stdio and HTTP transports, making it suitable for a variety of integration and debugging scenarios.

## Features
- **Stdio and HTTP transport modes**: Proxy MCP traffic over standard input/output or HTTP.
- **Flexible configuration**: Supports command-line flags, environment variables, and configuration files (YAML, JSON, TOML) via Viper.
- **Unique log file per run**: All traffic is logged to a uniquely named file for traceability.
- **Simple integration**: Easily drop into existing MCP workflows for debugging or auditing.

## Usage
- **Stdio mode**: Forwards stdin/stdout between client and sqlpp MCP server.
- **HTTP mode**: Listens on a local port and forwards HTTP requests to the sqlpp MCP server, logging all traffic.

## Configuration
- Command-line flags (with short and long forms)
- Environment variables (automatically picked up by Viper)
- Optional config file (YAML, JSON, TOML)

## Example
```
# Stdio mode
./mcp_sqlpp_proxy --transport stdio

# HTTP mode, listen on 8080, forward to sqlpp on 8811
./mcp_sqlpp_proxy --transport http --port 8080 --xfer-port 8811
```

## Logging
Each run creates a log file named `mcp_sqlpp_proxy_<pid>_<timestamp>.log` containing all proxied traffic for auditing and debugging.

## License
MIT
