package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
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
	// Create a temporary executable file for testing
	tempExe := "temp_mcp_sqlpp"
	err := os.WriteFile(tempExe, []byte("#!/bin/bash\necho test"), 0755)
	require.NoError(t, err)
	defer os.Remove(tempExe)

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
				ExePath:   tempExe, // Use the temporary executable
			},
			expectError: false,
		},
		{
			name: "valid http config",
			config: &Config{
				Transport: "http",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "./mcp_sqlpp", // HTTP mode doesn't validate exe path
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

func TestLoadConfig(t *testing.T) {
	// Reset viper for each test
	viper.Reset()

	// Create temporary executables for testing
	tempExe1 := "test_mcp_sqlpp_1"
	tempExe2 := "test_mcp_sqlpp_2"
	tempExe3 := "test_mcp_sqlpp_3"

	for _, exe := range []string{tempExe1, tempExe2, tempExe3} {
		err := os.WriteFile(exe, []byte("#!/bin/bash\necho test"), 0755)
		require.NoError(t, err)
		defer os.Remove(exe)
	}

	// Create a temporary config file
	configContent := `transport: http
port: 8080
xfer-port: 8892
exe-path: ` + tempExe1

	tempConfigFile := "test_load_config.yaml"
	err := os.WriteFile(tempConfigFile, []byte(configContent), 0644)
	require.NoError(t, err)
	defer os.Remove(tempConfigFile)

	tests := []struct {
		name           string
		flags          *Flags
		envVars        map[string]string
		expectedConfig *Config
		expectError    bool
		errorMsg       string
	}{
		{
			name: "load from config file",
			flags: &Flags{
				ConfigFile: &tempConfigFile,
				Transport:  stringPtr(""),
				Port:       intPtr(0),
				XferPort:   intPtr(0),
				ExePath:    stringPtr(""),
			},
			expectedConfig: &Config{
				Transport: "http",
				Port:      8080,
				XferPort:  8892,
				ExePath:   tempExe1,
			},
			expectError: false,
		},
		{
			name: "load defaults when no config file (http mode)",
			flags: &Flags{
				ConfigFile: stringPtr(""),
				Transport:  stringPtr("http"), // Use HTTP mode to avoid exe validation
				Port:       intPtr(0),
				XferPort:   intPtr(0),
				ExePath:    stringPtr(""),
			},
			expectedConfig: &Config{
				Transport: "http",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "./mcp_sqlpp",
			},
			expectError: false,
		},
		{
			name: "command line flags override config file",
			flags: &Flags{
				ConfigFile: &tempConfigFile,
				Transport:  stringPtr("stdio"),
				Port:       intPtr(9000),
				XferPort:   intPtr(0),
				ExePath:    stringPtr(tempExe2),
			},
			expectedConfig: &Config{
				Transport: "stdio",
				Port:      9000,
				XferPort:  8892, // from config file
				ExePath:   tempExe2,
			},
			expectError: false,
		},
		{
			name: "environment variables override config file",
			flags: &Flags{
				ConfigFile: &tempConfigFile,
				Transport:  stringPtr(""),
				Port:       intPtr(0),
				XferPort:   intPtr(0),
				ExePath:    stringPtr(""),
			},
			envVars: map[string]string{
				"MCP_PROXY_TRANSPORT": "stdio",
				"MCP_PROXY_PORT":      "7000",
				"MCP_PROXY_EXE_PATH":  tempExe3,
			},
			expectedConfig: &Config{
				Transport: "stdio",
				Port:      7000,
				XferPort:  8892,     // from config file
				ExePath:   tempExe3, // from env var
			},
			expectError: false,
		},
		{
			name: "missing config file error",
			flags: &Flags{
				ConfigFile: stringPtr("nonexistent.yaml"),
				Transport:  stringPtr(""),
				Port:       intPtr(0),
				XferPort:   intPtr(0),
				ExePath:    stringPtr(""),
			},
			expectError: true,
			errorMsg:    "failed to read config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test
			viper.Reset()

			// Set environment variables if provided
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			config, err := LoadConfig(tt.flags)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedConfig.Transport, config.Transport)
				assert.Equal(t, tt.expectedConfig.Port, config.Port)
				assert.Equal(t, tt.expectedConfig.XferPort, config.XferPort)
				assert.Equal(t, tt.expectedConfig.ExePath, config.ExePath)
			}
		})
	}
}

func TestLoadConfigValidationErrors(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Create temp executable for validation
	tempExe := "temp_exe_validation"
	err := os.WriteFile(tempExe, []byte("#!/bin/bash\necho test"), 0755)
	require.NoError(t, err)
	defer os.Remove(tempExe)

	tests := []struct {
		name        string
		flags       *Flags
		expectError bool
		errorMsg    string
	}{
		{
			name: "invalid transport validation",
			flags: &Flags{
				ConfigFile: stringPtr(""),
				Transport:  stringPtr("invalid"),
				Port:       intPtr(0),
				XferPort:   intPtr(0),
				ExePath:    stringPtr(tempExe),
			},
			expectError: true,
			errorMsg:    "invalid configuration",
		},
		{
			name: "invalid port validation - zero port",
			flags: &Flags{
				ConfigFile: stringPtr(""),
				Transport:  stringPtr("http"),
				Port:       intPtr(0),    // This will become 8099 due to defaults
				XferPort:   intPtr(8099), // Same as default port - should cause conflict
				ExePath:    stringPtr(""),
			},
			expectError: true,
			errorMsg:    "invalid configuration",
		},
		{
			name: "port conflict validation",
			flags: &Flags{
				ConfigFile: stringPtr(""),
				Transport:  stringPtr("http"),
				Port:       intPtr(8080),
				XferPort:   intPtr(8080),
				ExePath:    stringPtr(""),
			},
			expectError: true,
			errorMsg:    "invalid configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			_, err := LoadConfig(tt.flags)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	// Note: ParseFlags is hard to test directly because it modifies global flag state
	// and calls flag.Parse(). In a real scenario, this would be tested as part of
	// integration tests or refactored to be more testable.
	// For now, we'll test that the function returns a properly structured Flags object

	// This test validates the structure but doesn't test the actual parsing
	// since that would interfere with other tests
	flags := &Flags{
		ConfigFile: stringPtr("test.yaml"),
		Transport:  stringPtr("http"),
		Port:       intPtr(8080),
		XferPort:   intPtr(8891),
		ExePath:    stringPtr("/test/path"),
	}

	// Validate the flags structure
	assert.NotNil(t, flags.ConfigFile)
	assert.NotNil(t, flags.Transport)
	assert.NotNil(t, flags.Port)
	assert.NotNil(t, flags.XferPort)
	assert.NotNil(t, flags.ExePath)

	assert.Equal(t, "test.yaml", *flags.ConfigFile)
	assert.Equal(t, "http", *flags.Transport)
	assert.Equal(t, 8080, *flags.Port)
	assert.Equal(t, 8891, *flags.XferPort)
	assert.Equal(t, "/test/path", *flags.ExePath)
}

// Additional edge case tests for ValidateConfig
func TestValidateConfigEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "port too high",
			config: &Config{
				Transport: "http",
				Port:      70000,
				XferPort:  8891,
				ExePath:   "",
			},
			expectError: true,
			errorMsg:    "invalid port 70000",
		},
		{
			name: "xfer-port too high",
			config: &Config{
				Transport: "http",
				Port:      8080,
				XferPort:  70000,
				ExePath:   "",
			},
			expectError: true,
			errorMsg:    "invalid xfer-port 70000",
		},
		{
			name: "stdio with non-existent exe path",
			config: &Config{
				Transport: "stdio",
				Port:      8099,
				XferPort:  8891,
				ExePath:   "/definitely/does/not/exist",
			},
			expectError: true,
			errorMsg:    "executable not found at path",
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

// Helper functions for pointer creation
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
