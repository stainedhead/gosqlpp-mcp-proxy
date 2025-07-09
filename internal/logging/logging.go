package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger represents a structured logger for the MCP SQLPP Proxy
type Logger struct {
	*log.Logger
	filePath string
	file     *os.File
}

// LogConfig holds configuration for logging
type LogConfig struct {
	FilePath   string // Custom file path (optional)
	Prefix     string // Log prefix (optional)
	UseStdout  bool   // Also log to stdout
	TimeFormat string // Custom time format (optional)
}

// New creates a new logger instance with a unique log file
func New(config *LogConfig) (*Logger, error) {
	var filePath string
	var prefix string

	if config != nil && config.FilePath != "" {
		filePath = config.FilePath
	} else {
		// Generate unique log file name using timestamp and PID
		filePath = fmt.Sprintf("mcp_sqlpp_proxy_%d_%d.log", os.Getpid(), time.Now().UnixNano())
	}

	if config != nil && config.Prefix != "" {
		prefix = config.Prefix
	}

	// Create/open log file
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file '%s': %w", filePath, err)
	}

	var flags int = log.LstdFlags
	if config != nil && config.TimeFormat != "" {
		// Custom time format would require more complex implementation
		// For now, we'll use standard flags
		flags = log.LstdFlags
	}

	logger := &Logger{
		Logger:   log.New(logFile, prefix, flags),
		filePath: filePath,
		file:     logFile,
	}

	// Also log to stdout if requested
	if config != nil && config.UseStdout {
		// TODO: Implement multi-writer for both file and stdout
		// For now, we'll just use file logging
	}

	return logger, nil
}

// NewDefault creates a logger with default settings
func NewDefault() (*Logger, error) {
	return New(nil)
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		err := l.file.Close()
		l.file = nil // Set to nil to prevent double close
		return err
	}
	return nil
}

// GetFilePath returns the path to the log file
func (l *Logger) GetFilePath() string {
	return l.filePath
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	l.Printf("[INFO] %s", msg)
}

// Infof logs an informational message with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Printf("[INFO] "+format, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.Printf("[ERROR] %s", msg)
}

// Errorf logs an error message with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Printf("[ERROR] "+format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.Printf("[DEBUG] %s", msg)
}

// Debugf logs a debug message with formatting
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Printf("[DEBUG] "+format, args...)
}

// TrafficIn logs incoming traffic
func (l *Logger) TrafficIn(msg string) {
	l.Printf("[IN] %s", msg)
}

// TrafficOut logs outgoing traffic
func (l *Logger) TrafficOut(msg string) {
	l.Printf("[OUT] %s", msg)
}

// HTTPIn logs incoming HTTP request
func (l *Logger) HTTPIn(method, url string) {
	l.Printf("[HTTP IN] %s %s", method, url)
}

// HTTPInBody logs incoming HTTP request body
func (l *Logger) HTTPInBody(body string) {
	l.Printf("[HTTP IN BODY] %s", body)
}

// HTTPOut logs outgoing HTTP response
func (l *Logger) HTTPOut(statusCode int, body string) {
	l.Printf("[HTTP OUT] %d %s", statusCode, body)
}

// HTTPError logs HTTP-related errors
func (l *Logger) HTTPError(err error) {
	l.Printf("[HTTP ERROR] %v", err)
}

// Startup logs application startup information
func (l *Logger) Startup(msg string) {
	l.Printf("[STARTUP] %s", msg)
}

// Startupf logs application startup information with formatting
func (l *Logger) Startupf(format string, args ...interface{}) {
	l.Printf("[STARTUP] "+format, args...)
}

// Fatal logs a fatal error and exits the application
func (l *Logger) Fatal(msg string) {
	l.Printf("[FATAL] %s", msg)
	if l.file != nil {
		l.file.Close()
	}
	os.Exit(1)
}

// Fatalf logs a fatal error with formatting and exits the application
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Printf("[FATAL] "+format, args...)
	if l.file != nil {
		l.file.Close()
	}
	os.Exit(1)
}
