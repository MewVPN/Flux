package main

import (
	"log"
	"time"

	"flux/internal/config"
	"flux/internal/http"
	"flux/internal/wg"
)

var startTime = time.Now()

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := wg.Ensure(cfg); err != nil {
		log.Fatal(err)
	}

	router := http.NewRouter(cfg, startTime)

	log.Println("agent listening on :8080")
	router.Run(":8080")
}
