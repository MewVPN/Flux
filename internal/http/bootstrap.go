package http

import (
	"flux/internal/config"
	"flux/internal/util"

	"github.com/gin-gonic/gin"
)

// POST /bootstrap
func Bootstrap(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.SharedSecret != "" {
			c.JSON(409, gin.H{
				"error": "already bootstrapped",
			})
			return
		}

		cfg.UUID = util.UUIDv4()
		cfg.Name = util.AgentName()

		cfg.SharedSecret = util.Secret(32)

		if err := config.Save(cfg); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"uuid":         cfg.UUID,
			"name":         cfg.Name,
			"sharedSecret": cfg.SharedSecret,
			"message":      "agent successfully initialized",
		})
	}
}
