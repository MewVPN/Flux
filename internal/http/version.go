package http

import (
	"flux/internal/config"

	"github.com/gin-gonic/gin"
)

func VersionHandler(cfg *config.Config, version, commit, buildDate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":         cfg.Name,
			"version":      version,
			"commit":       commit,
			"buildDate":    buildDate,
			"bootstrapped": cfg.SharedSecret != "",
		})
	}
}
