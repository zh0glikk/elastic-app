package config

import (
	"encoding/json"
	"os"
)

type ListenerConfig struct {
	Addr string `json:"addr"`
}

type ElasticCfg struct {
	URL string `json:"url"`
}

type Config struct {
	ElasticCfg     ElasticCfg     `json:"elastic_cfg"`
	ListenerConfig ListenerConfig `json:"listener_config"`
}

func SetupConfig(filePath string, cfg *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return err
	}

	return nil
}
