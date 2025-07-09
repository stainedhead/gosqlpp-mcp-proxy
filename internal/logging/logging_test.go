package logging

import (
	"os"
	"strings"
	"testing"
)

func TestNewDefault(t *testing.T) {
	logger, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() failed: %v", err)
	}
	defer logger.Close()

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	if logger.GetFilePath() == "" {
		t.Fatal("Expected log file path to be set")
	}

	// Verify the log file was created
	if _, err := os.Stat(logger.GetFilePath()); os.IsNotExist(err) {
		t.Fatalf("Log file was not created: %s", logger.GetFilePath())
	}

	// Clean up
	defer os.Remove(logger.GetFilePath())
}

func TestNewWithConfig(t *testing.T) {
	config := &LogConfig{
		FilePath: "test_log.log",
		Prefix:   "[TEST] ",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer logger.Close()

	if logger.GetFilePath() != "test_log.log" {
		t.Errorf("Expected file path 'test_log.log', got '%s'", logger.GetFilePath())
	}

	// Clean up
	defer os.Remove(logger.GetFilePath())
}

func TestLoggerMethods(t *testing.T) {
	logger, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() failed: %v", err)
	}
	defer logger.Close()
	defer os.Remove(logger.GetFilePath())

	// Test various logging methods
	logger.Info("test info message")
	logger.Infof("test info message with format: %s", "test")
	logger.Error("test error message")
	logger.Errorf("test error message with format: %d", 123)
	logger.Debug("test debug message")
	logger.Debugf("test debug message with format: %v", true)
	logger.TrafficIn("test input traffic")
	logger.TrafficOut("test output traffic")
	logger.HTTPIn("GET", "/test")
	logger.HTTPInBody("test body")
	logger.HTTPOut(200, "OK")
	logger.HTTPError(err)
	logger.Startup("test startup message")
	logger.Startupf("test startup message with format: %s", "formatted")

	// Read the log file content
	content, err := os.ReadFile(logger.GetFilePath())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify that different log types are present
	expectedPrefixes := []string{
		"[INFO]",
		"[ERROR]",
		"[DEBUG]",
		"[IN]",
		"[OUT]",
		"[HTTP IN]",
		"[HTTP IN BODY]",
		"[HTTP OUT]",
		"[HTTP ERROR]",
		"[STARTUP]",
	}

	for _, prefix := range expectedPrefixes {
		if !strings.Contains(logContent, prefix) {
			t.Errorf("Expected log to contain '%s', but it didn't. Log content:\n%s", prefix, logContent)
		}
	}
}

func TestLoggerClose(t *testing.T) {
	logger, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() failed: %v", err)
	}
	defer os.Remove(logger.GetFilePath())

	// Close the logger
	err = logger.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Try to close again (should not fail)
	err = logger.Close()
	if err != nil {
		t.Errorf("Second Close() call failed: %v", err)
	}
}

func TestNewWithInvalidPath(t *testing.T) {
	config := &LogConfig{
		FilePath: "/nonexistent/directory/test.log",
	}

	_, err := New(config)
	if err == nil {
		t.Fatal("Expected New() to fail with invalid path, but it succeeded")
	}
}

func TestFatalBehavior(t *testing.T) {
	// This test is tricky because Fatal calls os.Exit(1)
	// We can't easily test this without a subprocess, so we'll skip this for now
	// In a more comprehensive test suite, we would use a subprocess or mock os.Exit
	t.Skip("Fatal behavior testing requires subprocess or mocking os.Exit")
}
