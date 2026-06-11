package main

import (
	"fmt"
	"kiro-go/config"
	"kiro-go/contentfilter"
	"kiro-go/logger"
	"kiro-go/pool"
	"kiro-go/proxy"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	configPath := "data/config.json"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	if err := config.Init(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(config.GetLogLevel())

	contentfilter.Load("context-filtes/filters.json")
	if logger.GetLevel() == logger.LevelDebug {
		contentfilter.SetAuditMode(true)
	}

	if envPassword := os.Getenv("ADMIN_PASSWORD"); envPassword != "" {
		config.SetPassword(envPassword)
	}

	pool.GetPool()

	handler := proxy.NewHandler()

	host := config.GetHost()
	port := config.GetPort()
	addr := fmt.Sprintf("%s:%d", host, port)
	displayHost := host
	if displayHost == "0.0.0.0" || displayHost == "" {
		displayHost = "127.0.0.1"
	}
	displayAddr := fmt.Sprintf("%s:%d", displayHost, port)
	logger.Infof("Kiro-Cybxai starting on http://%s (log level: %s)", addr, logger.LevelName(logger.GetLevel()))
	logger.Infof("Dashboard: http://%s", displayAddr)
	logger.Infof("Claude API: http://%s/v1/messages", displayAddr)
	logger.Infof("OpenAI API: http://%s/v1/chat/completions", displayAddr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
