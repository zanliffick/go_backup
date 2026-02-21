package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerIP   string `json:"server_ip"`
	Username   string `json:"username"`
	LocalDir   string `json:"local_dir"`
	RemoteDir  string `json:"remote_dir"`
	FilePrefix string `json:"file_prefix"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.ServerIP == "" {
		return nil, fmt.Errorf("config missing required field: server_ip")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("config missing required field: username")
	}
	if cfg.LocalDir == "" {
		return nil, fmt.Errorf("config missing required field: local_dir")
	}
	if cfg.RemoteDir == "" {
		return nil, fmt.Errorf("config missing required field: remote_dir")
	}
	if cfg.FilePrefix == "" {
		return nil, fmt.Errorf("config missing required field: file_prefix")
	}

	return &cfg, nil
}
