package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Registry struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Registries []Registry `yaml:"registries"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "harbor.yaml")
}

func EnsureDefaultConfig() {
	path := configPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0755)
		defaultConfig := []byte(`registries:
  - name: dockerhub
    url: https://hub.docker.com`)
		os.WriteFile(path, defaultConfig, 0644)
	}
}

func LoadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}
