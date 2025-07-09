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
   "github.com/spf13/viper"
   flag "github.com/spf13/pflag"
)


func main() {
	// Add config file flag
   configFile := flag.String("config", "", "Path to config file (yaml/json/toml)")
   transport := flag.StringP("transport", "t", "stdio", "Transport mode: stdio or http")
   port := flag.IntP("port", "p", 8099, "Port to listen on (HTTP mode)")
   xferPort := flag.IntP("xfer-port", "x", 8891, "Port where mcp_sqlpp is running (HTTP mode)")
   flag.Parse()

   viper.SetDefault("transport", "stdio")
   viper.SetDefault("port", 8099)
   viper.SetDefault("xfer-port", 8891)

   viper.BindEnv("transport")
   viper.BindEnv("port")
   viper.BindEnv("xfer-port")

   viper.BindPFlag("transport", flag.Lookup("transport"))
   viper.BindPFlag("port", flag.Lookup("port"))
   viper.BindPFlag("xfer-port", flag.Lookup("xfer-port"))

   // Use flag values
   viper.Set("transport", *transport)
   viper.Set("port", *port)
   viper.Set("xfer-port", *xferPort)

	// Load config file if provided
	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}
	viper.ReadInConfig() // Ignore error if no config file

	transportVal := viper.GetString("transport")
	portVal := viper.GetInt("port")
	xferPortVal := viper.GetInt("xfer-port")

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
		logger.Println("Starting in stdio mode")
		runStdioProxy(logger)
	case "http":
		logger.Printf("Starting in http mode on port %d, forwarding to localhost:%d", portVal, xferPortVal)
		runHTTPProxy(portVal, xferPortVal, logger)
	default:
		logger.Fatalf("Unknown transport: %s", transportVal)
	}
}

func runStdioProxy(logger *log.Logger) {
	cmd := exec.Command("/Users/mma0975/.sqlpp/mcp_sqlpp", "-t", "stdio")
	mcpIn, _ := cmd.StdinPipe()
	mcpOut, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Fatalf("Failed to start mcp_sqlpp: %v", err)
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
