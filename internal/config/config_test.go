package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "stdio", config.Transport)
	assert.Equal(t, 8099, config.Port)
	assert.Equal(t, 8891, config.XferPort)
	assert.Equal(t, "./mcp_sqlpp", config.ExePath)
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid stdio config",
			config: &Config{
				Transport: "stdio",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "./mcp_sqlpp",
			},
			expectError: false,
		},
		{
			name: "valid http config",
			config: &Config{
				Transport: "http",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "./mcp_sqlpp",
			},
			expectError: false,
		},
		{
			name: "invalid transport",
			config: &Config{
				Transport: "invalid",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "./mcp_sqlpp",
			},
			expectError: true,
			errorMsg:    "invalid transport mode",
		},
		{
			name: "invalid port for http",
			config: &Config{
				Transport: "http",
				Port:      0,
				XferPort:  8891,
				ExePath:   "./mcp_sqlpp",
			},
			expectError: true,
			errorMsg:    "invalid port",
		},
		{
			name: "same port and xfer-port",
			config: &Config{
				Transport: "http",
				Port:      8099,
				XferPort:  8099,
				ExePath:   "./mcp_sqlpp",
			},
			expectError: true,
			errorMsg:    "port (8099) and xfer-port (8099) cannot be the same",
		},
		{
			name: "empty exe-path for stdio",
			config: &Config{
				Transport: "stdio",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "",
			},
			expectError: true,
			errorMsg:    "exe-path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenerateExampleConfig(t *testing.T) {
	tempFile := "test_config.yaml"
	defer os.Remove(tempFile)

	err := GenerateExampleConfig(tempFile)
	require.NoError(t, err)

	// Check that file was created
	_, err = os.Stat(tempFile)
	require.NoError(t, err)

	// Check file content
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "transport: stdio")
	assert.Contains(t, contentStr, "port: 8099")
	assert.Contains(t, contentStr, "xfer-port: 8891")
	assert.Contains(t, contentStr, "exe-path: ./mcp_sqlpp")
	assert.Contains(t, contentStr, "MCP SQLPP Proxy Configuration")
}

func TestConfigString(t *testing.T) {
	config := &Config{
		Transport: "http",
		Port:      8080,
		XferPort:  8891,
		ExePath:   "/usr/local/bin/mcp_sqlpp",
	}

	expected := "Config{Transport: http, Port: 8080, XferPort: 8891, ExePath: /usr/local/bin/mcp_sqlpp}"
	assert.Equal(t, expected, config.String())
}
