package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

const (
	ConfigDir  = ".chat-cli"
	ConfigFile = "config.yaml"
)

var (
	globalConfig *Config
	once         sync.Once
	loadErr      error
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
	Username      string `yaml:"username"`
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.ServerAddress) == "" {
		return fmt.Errorf("server address cannot be empty")
	}
	if strings.TrimSpace(c.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	return nil
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ConfigDir, ConfigFile), nil
}

func Save(cfg *Config, path string) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func Load(parts ...string) (*Config, error) {
	path := ""
	if len(parts) == 0 {
		defaultPath, err := GetConfigPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get config path: %w", err)
		}
		path = defaultPath
	} else {
		path = parts[0]
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

func Get(parts ...string) (*Config, error) {
	path := ""
	if len(parts) == 0 {
		defaultPath, err := GetConfigPath()
		if err != nil {
			return nil, err
		}
		path = defaultPath
	} else {
		path = parts[0]
	}
	once.Do(func() {
		globalConfig, loadErr = Load(path)
	})
	return globalConfig, loadErr
}

func Set(cfg *Config, parts ...string) error {
	path := ""
	if len(parts) == 0 {
		defaultPath, err := GetConfigPath()
		if err != nil {
			return err
		}
		path = defaultPath
	} else {
		path = parts[0]
	}
	if err := Save(cfg, path); err != nil {
		return err
	}
	globalConfig = cfg
	return nil
}

func Reset() {
	once = sync.Once{}
	globalConfig = nil
	loadErr = nil
}
