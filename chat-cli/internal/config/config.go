package internal

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	Username      string `mapstructure:"username"`
}

var Cfg Config

func LoadConfig(configName string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") 
	viper.AddConfigPath("./config")

	// Environment variable override
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults or env vars")
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return err
	}

	return nil
}

func Init() {
	if err := LoadConfig("config"); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
}