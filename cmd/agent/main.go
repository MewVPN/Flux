package main

import (
	"log"
	"os"
	"time"

	"flux/internal/config"
	httpapi "flux/internal/http"

	"github.com/gin-gonic/gin"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
	startTime = time.Now()
)

func main() {
	log.Printf("Starting Flux...")
	log.Printf("Version: %s | Commit: %s | BuildDate: %s", version, commit, buildDate)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := httpapi.NewRouter(cfg, version, commit, buildDate, startTime)

	log.Printf("Flux listening on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
