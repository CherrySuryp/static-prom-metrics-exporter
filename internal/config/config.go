package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type StaticMetric struct {
	Name  string `yaml:"name"`
	Value uint64 `yaml:"value"`
}

type Config struct {
	BasicAuth     map[string]string `yaml:"basic-auth"`
	StaticMetrics []StaticMetric    `yaml:"static-metrics"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	print()
	return &cfg
}
