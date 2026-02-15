package http

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ServiceResult struct {
	Status     string `json:"status"`
	LatencyMs  int64  `json:"latencyMs"`
	JitterMs   int64  `json:"jitterMs"`
	HTTPStatus int    `json:"httpStatus"`
	Samples    int    `json:"samples"`
}

type ServicesHealth struct {
	LastChecked   time.Time                `json:"lastChecked"`
	OverallStatus string                   `json:"overallStatus"`
	Services      map[string]ServiceResult `json:"services"`
}

var (
	healthMu sync.RWMutex
	cached   ServicesHealth
)

var targets = map[string]string{
	"chatgpt":   "https://chat.openai.com",
	"instagram": "https://www.instagram.com",
	"telegram":  "https://web.telegram.org",
	"whatsapp":  "https://www.whatsapp.com",
}

func StartServicesHealthChecker() {
	go func() {
		runCheck()

		ticker := time.NewTicker(2 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			runCheck()
		}
	}()
}

func runCheck() {
	results := make(map[string]ServiceResult)
	upCount := 0

	for name, url := range targets {
		res := measureService(url)
		results[name] = res
		if res.Status == "up" {
			upCount++
		}
	}

	overall := "critical"
	if upCount == len(targets) {
		overall = "healthy"
	} else if upCount > 0 {
		overall = "degraded"
	}

	healthMu.Lock()
	cached = ServicesHealth{
		LastChecked:   time.Now().UTC(),
		OverallStatus: overall,
		Services:      results,
	}
	healthMu.Unlock()
}

func HealthServices(c *gin.Context) {
	healthMu.RLock()
	defer healthMu.RUnlock()
	c.JSON(http.StatusOK, cached)
}

func measureService(url string) ServiceResult {
	const samples = 3
	var latencies []int64
	var lastStatus int

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	for i := 0; i < samples; i++ {
		start := time.Now()
		resp, err := client.Get(url)
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			return ServiceResult{
				Status:    "down",
				LatencyMs: 0,
				JitterMs:  0,
				Samples:   i,
			}
		}

		lastStatus = resp.StatusCode
		resp.Body.Close()

		latencies = append(latencies, elapsed)
		time.Sleep(100 * time.Millisecond)
	}

	return ServiceResult{
		Status:     "up",
		LatencyMs:  average(latencies),
		JitterMs:   calculateJitter(latencies),
		HTTPStatus: lastStatus,
		Samples:    samples,
	}
}

func average(values []int64) int64 {
	var sum int64
	for _, v := range values {
		sum += v
	}
	return sum / int64(len(values))
}

func calculateJitter(values []int64) int64 {
	if len(values) < 2 {
		return 0
	}

	var total int64
	for i := 1; i < len(values); i++ {
		diff := values[i] - values[i-1]
		if diff < 0 {
			diff = -diff
		}
		total += diff
	}

	return total / int64(len(values)-1)
}
