package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds every server-tunable parameter.
type Config struct {
	ListenAddress   string        `mapstructure:"listen_address"`    // ":9000"
	Transport       string        `mapstructure:"transport"`         // "tcp" or "websocket"
	MaxMessageBytes int           `mapstructure:"max_message_bytes"` // 4096
	LogLevel        string        `mapstructure:"log_level"`         // "info", "debug", etc.
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`     // "5s"
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`      // "30s"
}

var Cfg Config // populated by Load()

// Load reads config from file + env vars.
// Call this once at startup (e.g., in cmd/server/main.go).
func Load(configName string) error {
	viper.SetConfigName(configName) // "config" → config.yaml
	viper.SetConfigType("yaml")

	// Locations to search—in order.
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.chat-server")
	viper.AddConfigPath("./config")

	// Allow env vars like CHAT_SERVER_LISTEN_ADDRESS to override.
	viper.SetEnvPrefix("chat_server")
	viper.AutomaticEnv()

	// Defaults keep the server runnable even without a file.
	viper.SetDefault("listen_address", ":9000")
	viper.SetDefault("transport", "tcp")
	viper.SetDefault("max_message_bytes", 4096)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("write_timeout", "5s")
	viper.SetDefault("read_timeout", "30s")

	// It’s okay if the file doesn’t exist; use defaults + env.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config: %w", err)
		}
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}
