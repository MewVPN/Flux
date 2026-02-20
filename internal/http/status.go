package http

import (
	"time"

	"flux/internal/config"
	"flux/internal/wg"

	"github.com/gin-gonic/gin"
)

func StatusHandler(cfg *config.Config, start time.Time) gin.HandlerFunc {
	return func(c *gin.Context) {

		uptime := int(time.Since(start).Seconds())
		bootstrapped := cfg.SharedSecret != ""

		wgRunning := wg.Running()
		wgPeers := wg.PeerCount()

		overall := "healthy"
		if !wgRunning {
			overall = "degraded"
		}

		c.JSON(200, gin.H{
			"unique_id":     cfg.UUID,
			"name":          cfg.Name,
			"bootstrapped":  bootstrapped,
			"uptimeSeconds": uptime,
			"overallStatus": overall,
			"wg": gin.H{
				"running": wgRunning,
				"peers":   wgPeers,
			},
		})
	}
}
