package http

import (
	"net/http"

	"flux/internal/config"
	"flux/internal/util"

	"github.com/gin-gonic/gin"
)

type BootstrapRequest struct {
	BootstrapHost string `json:"bootstrapHost" binding:"required"`
}

func Bootstrap(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		if cfg.SharedSecret != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "already bootstrapped",
			})
			return
		}

		var req BootstrapRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bootstrapHost is required",
			})
			return
		}

		cfg.UUID = util.UUIDv4()
		cfg.Name = util.AgentName()
		cfg.SharedSecret = util.Secret(32)

		cfg.BootstrapHost = req.BootstrapHost

		if err := config.Save(cfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"uuid":         cfg.UUID,
			"name":         cfg.Name,
			"sharedSecret": cfg.SharedSecret,
			"message":      "agent successfully initialized",
		})
	}
}
