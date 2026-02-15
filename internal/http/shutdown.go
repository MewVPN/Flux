package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"flux/internal/config"
)

func NotifyBootstrap(cfg *config.Config, startTime time.Time, version string) {
	if cfg.BootstrapHost == "" {
		return
	}

	payload := map[string]interface{}{
		"uuid":          cfg.UUID,
		"name":          cfg.Name,
		"uptimeSeconds": int(time.Since(startTime).Seconds()),
		"version":       version,
	}

	body, _ := json.Marshal(payload)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	url := cfg.BootstrapHost + "/agent/shutdown"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		log.Printf("Failed to notify bootstrap: %v", err)
	}
}
