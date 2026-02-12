package http

import (
	"flux/internal/config"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	serviceTimeout = 2 * time.Second
	checkInterval  = 1 * time.Hour
)

type serviceResult struct {
	Status     string `json:"status"`
	LatencyMs  int64  `json:"latencyMs,omitempty"`
	HTTPStatus int    `json:"httpStatus,omitempty"`
	Error      string `json:"error,omitempty"`
}

type servicesCache struct {
	sync.RWMutex
	LastChecked time.Time                `json:"lastChecked"`
	Services    map[string]serviceResult `json:"services"`
}

var (
	httpClient = &http.Client{Timeout: serviceTimeout}

	services = map[string]string{
		"instagram": "https://www.instagram.com/robots.txt",
		"whatsapp":  "https://web.whatsapp.com",
		"chatgpt":   "https://chat.openai.com",
		"telegram":  "https://api.telegram.org",
	}

	cache = servicesCache{
		Services: make(map[string]serviceResult),
	}
)

func HealthSelf(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",
			"bootstrapped": cfg.SharedSecret != "",
		})
	}
}

func HealthServices(c *gin.Context) {
	cache.RLock()
	defer cache.RUnlock()

	if cache.LastChecked.IsZero() {
		c.JSON(503, gin.H{
			"status":  "not_ready",
			"message": "services check not completed yet",
		})
		return
	}

	c.JSON(200, gin.H{
		"lastChecked": cache.LastChecked.UTC(),
		"services":    cache.Services,
	})
}

func StartServicesHealthChecker() {
	go func() {
		// first run immediately
		runServicesCheck()

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for range ticker.C {
			runServicesCheck()
		}
	}()
}

func runServicesCheck() {
	results := make(map[string]serviceResult)

	for name, url := range services {
		start := time.Now()

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			results[name] = serviceResult{
				Status: "error",
				Error:  err.Error(),
			}
			continue
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			results[name] = serviceResult{
				Status: "down",
				Error:  err.Error(),
			}
			continue
		}
		resp.Body.Close()

		results[name] = serviceResult{
			Status:     "up",
			LatencyMs:  time.Since(start).Milliseconds(),
			HTTPStatus: resp.StatusCode,
		}
	}

	cache.Lock()
	cache.Services = results
	cache.LastChecked = time.Now()
	cache.Unlock()
}
