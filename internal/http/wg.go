package http

import (
	"encoding/json"
	"io"

	"flux/internal/config"
	"flux/internal/wg"

	"github.com/gin-gonic/gin"
)

func WGList(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := wg.Request(cfg, "GET", "/client")
		if err != nil {
			c.JSON(502, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		// no filtering
		status := c.Query("status")
		if status == "" {
			c.Data(resp.StatusCode, "application/json", body)
			return
		}

		var clients []map[string]interface{}
		if err := json.Unmarshal(body, &clients); err != nil {
			c.JSON(500, gin.H{"error": "invalid wg-easy response"})
			return
		}

		filtered := make([]map[string]interface{}, 0)

		for _, c := range clients {
			enabled, _ := c["enabled"].(bool)

			if status == "active" && enabled {
				filtered = append(filtered, c)
			}
			if status == "disabled" && !enabled {
				filtered = append(filtered, c)
			}
		}

		c.JSON(200, filtered)
	}
}

// WGCreate POST /wg/peers
func WGCreate(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := wg.Request(cfg, "POST", "/client")
		if err != nil {
			c.JSON(502, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", body)
	}
}

// WGEnable POST /wg/peers/:id/enable
func WGEnable(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		resp, err := wg.Request(cfg, "POST", "/client/"+id+"/enable")
		if err != nil {
			c.JSON(502, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", body)
	}
}

// WGDisable POST /wg/peers/:id/disable
func WGDisable(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		resp, err := wg.Request(cfg, "POST", "/client/"+id+"/disable")
		if err != nil {
			c.JSON(502, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", body)
	}
}

// WGConfig GET /wg/peers/:id/config
func WGConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		resp, err := wg.Request(cfg, "GET", "/client/"+id+"/configuration")
		if err != nil {
			c.JSON(502, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		c.Data(
			resp.StatusCode,
			"text/plain; charset=utf-8",
			body,
		)
	}
}
