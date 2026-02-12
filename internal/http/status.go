package http

import (
	"time"

	"flux/internal/config"

	"github.com/gin-gonic/gin"
)

func Status(cfg *config.Config, start time.Time) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"uniqueId":      cfg.UUID,
			"name":          cfg.Name,
			"uptimeSeconds": int(time.Since(start).Seconds()),
		})
	}
}
