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

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// Command line flags
	configFile := flag.String("config", "", "Path to config file (yaml/json/toml)")
	transport := flag.StringP("transport", "t", "", "Transport mode: stdio or http")
	port := flag.IntP("port", "p", 0, "Port to listen on (HTTP mode)")
	xferPort := flag.IntP("xfer-port", "x", 0, "Port where mcp_sqlpp is running (HTTP mode)")
	exePath := flag.StringP("exe-path", "e", "", "Path to the mcp_sqlpp executable")
	flag.Parse()

	// Set config defaults
	viper.SetDefault("transport", "stdio")
	viper.SetDefault("port", 8099)
	viper.SetDefault("xfer-port", 8891)
	viper.SetDefault("exe-path", "./mcp_sqlpp")

	// Bind environment variables
	viper.BindEnv("transport")
	viper.BindEnv("port")
	viper.BindEnv("xfer-port")
	viper.BindEnv("exe-path")

	// Load config file if provided
	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	// Read config file and handle errors
	if err := viper.ReadInConfig(); err != nil {
		if *configFile != "" {
			// If a specific config file was requested but not found, that's an error
			log.Fatalf("Failed to read config file '%s': %v", *configFile, err)
		}
		// If no specific config file was requested, it's okay if default config doesn't exist
	}

	// Override config values with command line flags (flags take precedence)
	if *transport != "" {
		viper.Set("transport", *transport)
	}
	if *port != 0 {
		viper.Set("port", *port)
	}
	if *xferPort != 0 {
		viper.Set("xfer-port", *xferPort)
	}
	if *exePath != "" {
		viper.Set("exe-path", *exePath)
	}

	// Get final configuration values
	transportVal := viper.GetString("transport")
	portVal := viper.GetInt("port")
	xferPortVal := viper.GetInt("xfer-port")
	exePathVal := viper.GetString("exe-path")

	// Create a unique log file for each run using timestamp and PID
	logFileName := fmt.Sprintf("mcp_sqlpp_proxy_%d_%d.log", os.Getpid(), time.Now().UnixNano())
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	switch transportVal {
	case "stdio":
		logger.Printf("Starting in stdio mode with exe-path: %s", exePathVal)
		runStdioProxy(exePathVal, logger)
	case "http":
		logger.Printf("Starting in http mode on port %d, forwarding to localhost:%d", portVal, xferPortVal)
		runHTTPProxy(portVal, xferPortVal, logger)
	default:
		logger.Fatalf("Unknown transport: %s", transportVal)
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
