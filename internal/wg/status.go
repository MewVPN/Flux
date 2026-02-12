package wg

import (
	"fmt"

	"flux/internal/util"

	"github.com/gin-gonic/gin"
)

func StatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"running": running(),
			"ui":      fmt.Sprintf("http://%s:51821", util.DetectPublicIP()),
		})
	}
}
