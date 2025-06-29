package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/rohitsh16/go-blog/backend/db"
)

type Config struct {
	DatabaseConfig *db.DatabaseConfig `yaml:"DatabaseConfig"`
}

// LoadConfig loads important config for system to run
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
