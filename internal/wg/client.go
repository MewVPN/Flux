package wg

import (
	"encoding/base64"
	"net/http"

	"flux/internal/config"
)

const api = "http://127.0.0.1:51821/api"

func Request(cfg *config.Config, method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, api+path, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(cfg.WGEasyUser + ":" + cfg.WGEasyPassword),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func running() bool {
	resp, err := http.Get(api + "/client")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200 || resp.StatusCode == 401
}
