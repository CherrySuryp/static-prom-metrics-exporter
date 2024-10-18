package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type StaticMetric struct {
	Name  string `yaml:"name"`
	Help  string `yaml:"help" default:""`
	Value uint64 `yaml:"value"`
}

type Server struct {
	Port      string            `yaml:"port"`
	TlsCrt    string            `yaml:"tls_crt" default:""`
	TlsKey    string            `yaml:"tls_key" default:""`
	BasicAuth map[string]string `yaml:"basic-auth"`
}

type Config struct {
	Server        Server         `yaml:"server"`
	StaticMetrics []StaticMetric `yaml:"static_metrics"`
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
	return &cfg
}
