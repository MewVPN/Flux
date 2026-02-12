package http

import (
	"flux/internal/config"

	"github.com/gin-gonic/gin"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.SharedSecret == "" {
			c.AbortWithStatusJSON(503, gin.H{"error": "not bootstrapped"})
			return
		}

		if c.GetHeader("X-Agent-Secret") != cfg.SharedSecret {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
