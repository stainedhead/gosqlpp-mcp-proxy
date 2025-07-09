package config

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Transport string `mapstructure:"transport" yaml:"transport" json:"transport" toml:"transport"`
	Port      int    `mapstructure:"port" yaml:"port" json:"port" toml:"port"`
	XferPort  int    `mapstructure:"xfer-port" yaml:"xfer-port" json:"xfer-port" toml:"xfer-port"`
	ExePath   string `mapstructure:"exe-path" yaml:"exe-path" json:"exe-path" toml:"exe-path"`
}

// Flags represents command-line flags
type Flags struct {
	ConfigFile *string
	Transport  *string
	Port       *int
	XferPort   *int
	ExePath    *string
}

// DefaultConfig returns a Config struct with default values
func DefaultConfig() *Config {
	return &Config{
		Transport: "stdio",
		Port:      8099,
		XferPort:  8891,
		ExePath:   "./mcp_sqlpp",
	}
}

// ParseFlags parses command-line flags and returns a Flags struct
func ParseFlags() *Flags {
	flags := &Flags{
		ConfigFile: flag.String("config", "", "Path to config file (yaml/json/toml)"),
		Transport:  flag.StringP("transport", "t", "", "Transport mode: stdio or http"),
		Port:       flag.IntP("port", "p", 0, "Port to listen on (HTTP mode)"),
		XferPort:   flag.IntP("xfer-port", "x", 0, "Port where mcp_sqlpp is running (HTTP mode)"),
		ExePath:    flag.StringP("exe-path", "e", "", "Path to the mcp_sqlpp executable"),
	}
	flag.Parse()
	return flags
}

// LoadConfig loads configuration from multiple sources with proper precedence:
// 1. Command-line flags (highest priority)
// 2. Environment variables
// 3. Configuration file
// 4. Default values (lowest priority)
func LoadConfig(flags *Flags) (*Config, error) {
	// Set default values
	defaults := DefaultConfig()
	viper.SetDefault("transport", defaults.Transport)
	viper.SetDefault("port", defaults.Port)
	viper.SetDefault("xfer-port", defaults.XferPort)
	viper.SetDefault("exe-path", defaults.ExePath)

	// Bind environment variables with automatic env var name mapping
	viper.SetEnvPrefix("MCP_PROXY")
	viper.AutomaticEnv()

	// Explicitly bind environment variables for kebab-case config keys
	viper.BindEnv("transport", "MCP_PROXY_TRANSPORT")
	viper.BindEnv("port", "MCP_PROXY_PORT")
	viper.BindEnv("xfer-port", "MCP_PROXY_XFER_PORT")
	viper.BindEnv("exe-path", "MCP_PROXY_EXE_PATH")

	// Load config file if provided
	if *flags.ConfigFile != "" {
		viper.SetConfigFile(*flags.ConfigFile)
	} else {
		// Look for default config files in current directory
		viper.SetConfigName("mcp_sqlpp_proxy")
		viper.SetConfigType("yaml") // Default type, but will try others
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/etc/mcp-proxy")
	}

	// Read config file and handle errors appropriately
	if err := viper.ReadInConfig(); err != nil {
		if *flags.ConfigFile != "" {
			// If a specific config file was requested but not found, that's an error
			return nil, fmt.Errorf("failed to read config file '%s': %w", *flags.ConfigFile, err)
		}
		// If no specific config file was requested, it's okay if default config doesn't exist
		// We'll just use defaults and environment variables
	}

	// Override config values with command-line flags (flags take highest precedence)
	if *flags.Transport != "" {
		viper.Set("transport", *flags.Transport)
	}
	if *flags.Port != 0 {
		viper.Set("port", *flags.Port)
	}
	if *flags.XferPort != 0 {
		viper.Set("xfer-port", *flags.XferPort)
	}
	if *flags.ExePath != "" {
		viper.Set("exe-path", *flags.ExePath)
	}

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// ValidateConfig validates the configuration values
func ValidateConfig(config *Config) error {
	// Validate transport mode
	if config.Transport != "stdio" && config.Transport != "http" {
		return fmt.Errorf("invalid transport mode '%s': must be 'stdio' or 'http'", config.Transport)
	}

	// Validate ports for HTTP mode
	if config.Transport == "http" {
		if config.Port <= 0 || config.Port > 65535 {
			return fmt.Errorf("invalid port %d: must be between 1 and 65535", config.Port)
		}
		if config.XferPort <= 0 || config.XferPort > 65535 {
			return fmt.Errorf("invalid xfer-port %d: must be between 1 and 65535", config.XferPort)
		}
		if config.Port == config.XferPort {
			return fmt.Errorf("port (%d) and xfer-port (%d) cannot be the same", config.Port, config.XferPort)
		}
	}

	// Validate executable path exists for stdio mode
	if config.Transport == "stdio" {
		if config.ExePath == "" {
			return fmt.Errorf("exe-path cannot be empty for stdio transport mode")
		}
		// Check if executable exists and is executable
		if _, err := os.Stat(config.ExePath); os.IsNotExist(err) {
			return fmt.Errorf("executable not found at path '%s'", config.ExePath)
		}
	}

	return nil
}

// GenerateExampleConfig creates an example configuration file with default values and comments
func GenerateExampleConfig(filename string) error {
	exampleContent := `# MCP SQLPP Proxy Configuration
# This file demonstrates all available configuration options with their default values.
# You can use YAML, JSON, or TOML format for configuration files.

# Transport mode: "stdio" or "http"
# - stdio: Communicates via standard input/output (good for command-line tools)
# - http: Acts as HTTP proxy (good for web applications and services)
transport: stdio

# Port to listen on when using HTTP transport mode
# Only used when transport is set to "http"
# Default: 8099
port: 8099

# Port where the target mcp_sqlpp server is running (HTTP mode only)
# This is where the proxy will forward HTTP requests
# Only used when transport is set to "http"  
# Default: 8891
xfer-port: 8891

# Path to the mcp_sqlpp executable
# Can be absolute path or relative to current working directory
# Examples:
#   - ./mcp_sqlpp (relative path, default)
#   - /usr/local/bin/mcp_sqlpp (absolute path)
#   - ../gosqlpp-mcp-server/gosqlpp-mcp-server (relative to another directory)
# Default: ./mcp_sqlpp
exe-path: ./mcp_sqlpp

# Environment Variable Overrides:
# All configuration options can also be set via environment variables:
# - MCP_PROXY_TRANSPORT=http
# - MCP_PROXY_PORT=8080
# - MCP_PROXY_XFER_PORT=8891
# - MCP_PROXY_EXE_PATH=/usr/local/bin/mcp_sqlpp

# Command-line flags take the highest precedence and will override 
# both environment variables and config file values.
`

	return os.WriteFile(filename, []byte(exampleContent), 0644)
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	return fmt.Sprintf("Config{Transport: %s, Port: %d, XferPort: %d, ExePath: %s}",
		c.Transport, c.Port, c.XferPort, c.ExePath)
}
