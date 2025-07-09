package main

import (
	"testing"
	"os/exec"
)

func TestMainBuilds(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "mcp_sqlpp_proxy", "main.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
}
