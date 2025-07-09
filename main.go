package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/your-org/go-sqlpp-mcp-proxy/internal/config"
)

func main() {
	// Parse command-line flags
	flags := config.ParseFlags()

	// Load configuration from all sources
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create a unique log file for each run using timestamp and PID
	logFileName := fmt.Sprintf("mcp_sqlpp_proxy_%d_%d.log", os.Getpid(), time.Now().UnixNano())
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("Starting MCP SQLPP Proxy with configuration: %s", cfg.String())

	switch cfg.Transport {
	case "stdio":
		logger.Printf("Starting in stdio mode with exe-path: %s", cfg.ExePath)
		runStdioProxy(cfg.ExePath, logger)
	case "http":
		logger.Printf("Starting in http mode on port %d, forwarding to localhost:%d", cfg.Port, cfg.XferPort)
		runHTTPProxy(cfg.Port, cfg.XferPort, logger)
	default:
		logger.Fatalf("Unknown transport: %s", cfg.Transport)
	}
}

func runStdioProxy(exePath string, logger *log.Logger) {
	cmd := exec.Command(exePath, "-t", "stdio")
	mcpIn, _ := cmd.StdinPipe()
	mcpOut, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Fatalf("Failed to start mcp_sqlpp at '%s': %v", exePath, err)
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			logger.Printf("[IN] %s", line)
			io.WriteString(mcpIn, line+"\n")
		}
	}()

	scanner := bufio.NewScanner(mcpOut)
	for scanner.Scan() {
		line := scanner.Text()
		logger.Printf("[OUT] %s", line)
		fmt.Println(line)
	}

	cmd.Wait()
}

func runHTTPProxy(listenPort, xferPort int, logger *log.Logger) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[HTTP IN] %s %s", r.Method, r.URL)
		// Read request body
		body, _ := io.ReadAll(r.Body)
		logger.Printf("[HTTP IN BODY] %s", string(body))

		// Forward to mcp_sqlpp HTTP server
		url := fmt.Sprintf("http://localhost:%d%s", xferPort, r.URL.Path)
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			logger.Printf("[HTTP ERROR] %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("[HTTP ERROR] %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		logger.Printf("[HTTP OUT] %d %s", resp.StatusCode, string(respBody))

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	})

	logger.Printf("Listening on http://localhost:%d", listenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
