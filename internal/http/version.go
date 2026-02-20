package http

import "github.com/gin-gonic/gin"

func VersionHandler(version, commit, buildDate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":   version,
			"commit":    commit,
			"buildDate": buildDate,
		})
	}
}
