package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gosqlpp-mcp-proxy/internal/config"
	"gosqlpp-mcp-proxy/internal/logging"
)

func main() {
	// Parse command-line flags
	flags := config.ParseFlags()

	// Load configuration from all sources
	cfg, err := config.LoadConfig(flags)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create a logger with default settings (unique log file)
	logger, err := logging.NewDefault()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Startupf("Starting MCP SQLPP Proxy with configuration: %s", cfg.String())

	switch cfg.Transport {
	case "stdio":
		logger.Infof("Starting in stdio mode with exe-path: %s", cfg.ExePath)
		runStdioProxy(cfg.ExePath, logger)
	case "http":
		logger.Infof("Starting in http mode on port %d, forwarding to localhost:%d", cfg.Port, cfg.XferPort)
		runHTTPProxy(cfg.Port, cfg.XferPort, logger)
	default:
		logger.Fatalf("Unknown transport: %s", cfg.Transport)
	}
}

func runStdioProxy(exePath string, logger *logging.Logger) {
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
			logger.TrafficIn(line)
			io.WriteString(mcpIn, line+"\n")
		}
	}()

	scanner := bufio.NewScanner(mcpOut)
	for scanner.Scan() {
		line := scanner.Text()
		logger.TrafficOut(line)
		fmt.Println(line)
	}

	cmd.Wait()
}

func runHTTPProxy(listenPort, xferPort int, logger *logging.Logger) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.HTTPIn(r.Method, r.URL.String())
		// Read request body
		body, _ := io.ReadAll(r.Body)
		logger.HTTPInBody(string(body))

		// Forward to mcp_sqlpp HTTP server
		url := fmt.Sprintf("http://localhost:%d%s", xferPort, r.URL.Path)
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			logger.HTTPError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.HTTPError(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		logger.HTTPOut(resp.StatusCode, string(respBody))

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	})

	logger.Infof("Listening on http://localhost:%d", listenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
