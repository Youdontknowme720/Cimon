package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Config struct {
	TokenGitlab string `yaml:"tokengitlab"`
	TokenGithub string `yaml:"tokenGithub"`
}

var (
	config     *Config
	configPath string
	once       sync.Once
)

func ResolveConfigFilePath() (string, error) {
	if configPath != "" {
		return configPath, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}
	var configDir string
	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(home, ".config", "cimon")
	case "windows":
		configDir = filepath.Join(home, "AppData", "Roaming", "cimon")
	default:
		configDir = filepath.Join(home, ".config", "cimon")
	}
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	configPath = filepath.Join(configDir, "config.yaml")
	return configPath, nil
}

func LoadConfig() (*Config, error) {
	path, err := ResolveConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()
	cfg := &Config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := ResolveConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file for writing: %w", err)
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return nil
}
