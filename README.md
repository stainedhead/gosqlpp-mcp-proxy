
# gosqlpp_mcp_proxy

A robust Go-based proxy utility for forwarding and logging Model Context Protocol (MCP) traffic between clients and the sqlpp MCP server. Supports both stdio and HTTP transports, with flexible configuration via command-line flags, environment variables, and config files. All traffic is logged to a unique file per run for traceability and debugging.

## Features
- Stdio and HTTP transport modes
- Command-line flags (short and long forms)
- Environment variable and config file support (via Viper)
- Unique log file per run
- Simple integration for debugging and auditing

## Usage

### Stdio Mode
```
./mcp_sqlpp_proxy --transport stdio
```

### HTTP Mode
```
./mcp_sqlpp_proxy --transport http --port 8080 --xfer-port 8811
```

### With Config File
```
./mcp_sqlpp_proxy --config config.yaml
```

## Configuration
- `--transport`, `-t`: Transport mode (`stdio` or `http`)
- `--port`, `-p`: Port to listen on (HTTP mode)
- `--xfer-port`, `-x`: Port where mcp_sqlpp is running (HTTP mode)
- `--config`: Path to config file (yaml/json/toml)

Environment variables and config files are automatically picked up by Viper.

## Logging
Each run creates a log file named `mcp_sqlpp_proxy_<pid>_<timestamp>.log` containing all proxied traffic.

## License
MIT

---

See [docs/product.md](docs/product.md) for a detailed product overview.
