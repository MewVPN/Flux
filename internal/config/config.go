package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	SharedSecret string `json:"sharedSecret"`

	WGEasyUser     string `json:"wgEasyUser"`
	WGEasyPassword string `json:"wgEasyPassword"`
}

const path = "./agent-config.json"

func Load() (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	return &cfg, json.Unmarshal(b, &cfg)
}

func Save(cfg *Config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0600)
}
