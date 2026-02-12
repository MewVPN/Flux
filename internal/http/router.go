package http

import (
	"time"

	"flux/internal/config"
	"flux/internal/wg"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, start time.Time) *gin.Engine {
	r := gin.Default()

	StartServicesHealthChecker()

	r.POST("/bootstrap", Bootstrap(cfg))
	r.GET("/health/self", HealthSelf(cfg))
	r.GET("/health/services", HealthServices)

	protected := r.Group("/")
	protected.Use(Auth(cfg))

	protected.GET("/status", Status(cfg, start))
	protected.GET("/wg/status", wg.StatusHandler())

	// wg peers
	protected.GET("/wg/peers", WGList(cfg))
	protected.POST("/wg/peers", WGCreate(cfg))

	protected.POST("/wg/peers/:id/enable", WGEnable(cfg))
	protected.POST("/wg/peers/:id/disable", WGDisable(cfg))

	protected.GET("/wg/peers/:id/config", WGConfig(cfg))

	return r
}
