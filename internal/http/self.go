package http

import (
	"time"

	"flux/internal/config"

	"github.com/gin-gonic/gin"
)

func HealthSelf(cfg *config.Config, startTime time.Time) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"uuid":          cfg.UUID,
			"name":          cfg.Name,
			"bootstrapped":  cfg.SharedSecret != "",
			"uptimeSeconds": int(time.Since(startTime).Seconds()),
		})
	}
}
