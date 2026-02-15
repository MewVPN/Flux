package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := httpapi.NewRouter(cfg, version, commit, buildDate, startTime)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Flux listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	httpapi.NotifyBootstrap(cfg, startTime, version)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	log.Println("Flux stopped")
}
