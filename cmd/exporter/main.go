package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/l13t/github_rate_limit_exporter/internal/collector"
	"github.com/l13t/github_rate_limit_exporter/internal/config"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to configuration file (supports .yaml, .yml, .toml, .hcl)")
	version    = "dev"
)

func main() {
	flag.Parse()

	log.Printf("GitHub Rate Limit Exporter version %s", version)

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Loaded configuration with %d users", len(cfg.Users))
	log.Printf("Listen address: %s", cfg.ListenAddr)
	log.Printf("Metrics path: %s", cfg.MetricsPath)
	log.Printf("Poll interval: %d seconds", cfg.PollInterval)

	// Create collector
	c := collector.NewCollector(cfg.Users)

	// Register collector with Prometheus
	prometheus.MustRegister(c)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background polling
	go c.StartPolling(ctx, time.Duration(cfg.PollInterval)*time.Second)

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle(cfg.MetricsPath, promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>
<head><title>GitHub Rate Limit Exporter</title></head>
<body>
<h1>GitHub Rate Limit Exporter</h1>
<p><a href="` + cfg.MetricsPath + `">Metrics</a></p>
</body>
</html>`))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", cfg.ListenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")

	// Cancel background polling
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Exporter stopped")
}
