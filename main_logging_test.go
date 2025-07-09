package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"gosqlpp-mcp-proxy/internal/logging"
)

func TestLoggingIntegration(t *testing.T) {
	// Create a logger similar to how main.go does it
	logger, err := logging.NewDefault()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	logFile := logger.GetFilePath()
	defer os.Remove(logFile)

	// Test that logging methods work as expected in the context of main
	logger.Startupf("Starting MCP SQLPP Proxy with configuration: %s", "test-config")
	logger.Infof("Starting in stdio mode with exe-path: %s", "/test/path")
	logger.TrafficIn("test input")
	logger.TrafficOut("test output")
	logger.HTTPIn("GET", "/test")
	logger.HTTPInBody("test body")
	logger.HTTPOut(200, "OK")

	// Give a moment for the writes to complete
	time.Sleep(10 * time.Millisecond)

	// Read the log file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify expected log entries are present
	expectedEntries := []string{
		"[STARTUP] Starting MCP SQLPP Proxy with configuration: test-config",
		"[INFO] Starting in stdio mode with exe-path: /test/path",
		"[IN] test input",
		"[OUT] test output",
		"[HTTP IN] GET /test",
		"[HTTP IN BODY] test body",
		"[HTTP OUT] 200 OK",
	}

	for _, expected := range expectedEntries {
		if !strings.Contains(logContent, expected) {
			t.Errorf("Expected log to contain '%s', but it didn't. Log content:\n%s", expected, logContent)
		}
	}
}

func TestLoggerCleanupInMain(t *testing.T) {
	// This test ensures that logger cleanup (defer logger.Close()) works correctly
	// We simulate what happens in main.go

	var logFile string

	func() {
		logger, err := logging.NewDefault()
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close() // This is what main.go does

		logFile = logger.GetFilePath()
		logger.Info("test message")
	}() // Logger should be closed when this function exits

	// Verify the log file exists and contains our message
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file after cleanup: %v", err)
	}

	if !strings.Contains(string(content), "[INFO] test message") {
		t.Errorf("Log file doesn't contain expected message after cleanup")
	}

	// Clean up
	defer os.Remove(logFile)
}
