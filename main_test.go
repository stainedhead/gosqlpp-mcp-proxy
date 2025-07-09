package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestMainBuilds(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "mcp_sqlpp_proxy_test", "main.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	// Clean up
	defer os.Remove("mcp_sqlpp_proxy_test")
}

func TestMainHelpFlag(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "mcp_sqlpp_proxy_test", "main.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer os.Remove("mcp_sqlpp_proxy_test")

	// Test help flag
	cmd = exec.Command("./mcp_sqlpp_proxy_test", "--help")
	output, err := cmd.CombinedOutput()

	// Help flag should cause a non-zero exit (pflag behavior)
	if err == nil {
		t.Log("Help command completed successfully")
	}

	outputStr := string(output)
	expectedStrings := []string{
		"Usage of",
		"--config",
		"--transport",
		"--port",
		"--xfer-port",
		"--exe-path",
	}

	for _, expected := range expectedStrings {
		if !contains(outputStr, expected) {
			t.Errorf("Help output missing expected string: %s\nFull output:\n%s", expected, outputStr)
		}
	}
}

func TestMainConfigError(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "mcp_sqlpp_proxy_test", "main.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer os.Remove("mcp_sqlpp_proxy_test")

	// Test with non-existent config file
	cmd = exec.Command("./mcp_sqlpp_proxy_test", "--config", "nonexistent.yaml")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for non-existent config file")
	}

	outputStr := string(output)
	if !contains(outputStr, "Configuration error") {
		t.Errorf("Expected configuration error message, got: %s", outputStr)
	}
}

func TestMainValidationError(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "mcp_sqlpp_proxy_test", "main.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer os.Remove("mcp_sqlpp_proxy_test")

	// Test with invalid transport mode
	cmd = exec.Command("./mcp_sqlpp_proxy_test", "--transport", "invalid")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for invalid transport mode")
	}

	outputStr := string(output)
	if !contains(outputStr, "Configuration error") && !contains(outputStr, "invalid transport") {
		t.Errorf("Expected invalid transport error message, got: %s", outputStr)
	}
}

func TestMainWithTimeout(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "mcp_sqlpp_proxy_test", "main.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer os.Remove("mcp_sqlpp_proxy_test")

	// Create a fake executable for testing
	fakeExe := "fake_mcp_sqlpp"
	fakeScript := `#!/bin/bash
# Simulate a process that runs for a bit
sleep 1
echo "fake mcp server"
`
	err := os.WriteFile(fakeExe, []byte(fakeScript), 0755)
	if err != nil {
		t.Fatalf("Failed to create fake executable: %v", err)
	}
	defer os.Remove(fakeExe)

	// Test stdio mode with timeout
	cmd = exec.Command("./mcp_sqlpp_proxy_test", "--transport", "stdio", "--exe-path", fakeExe)

	// Start the command
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Kill after a short timeout
	time.AfterFunc(2*time.Second, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})

	// Wait for completion
	err = cmd.Wait()

	// It's expected that the command was killed or completed
	if err != nil {
		t.Logf("Command terminated as expected: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		(len(str) > len(substr) && containsSubstring(str, substr)))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
