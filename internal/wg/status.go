package wg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"flux/internal/util"

	"github.com/gin-gonic/gin"
)

const (
	containerName = "wg-easy"
	apiURL        = "http://127.0.0.1:51821/api/clients"
)

func StatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		running := Running()
		peers := 0

		if running {
			peers = PeerCount()
		}

		c.JSON(200, gin.H{
			"running": running,
			"peers":   peers,
			"ui":      fmt.Sprintf("http://%s:51821", util.DetectPublicIP()),
		})
	}
}

func Running() bool {
	out, err := exec.Command(
		"docker", "ps",
		"--filter", "name="+containerName,
		"--filter", "status=running",
		"--format", "{{.Names}}",
	).Output()

	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == containerName
}

func PeerCount() int {
	if !Running() {
		return 0
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var clients []interface{}

	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return 0
	}

	return len(clients)
}
