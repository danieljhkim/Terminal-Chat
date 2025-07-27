package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type InitConfig struct {
	ServerAddress string `yaml:"server_address"`
	Username      string `yaml:"username"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure the chat CLI",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter server address (e.g. localhost:9000): ")
		serverAddr, _ := reader.ReadString('\n')
		serverAddr = sanitizeInput(serverAddr)

		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = sanitizeInput(username)

		cfg := InitConfig{
			ServerAddress: serverAddr,
			Username:      username,
		}

		savePath := filepath.Join(os.Getenv("HOME"), ".chat-cli", "config.yaml")
		if err := saveConfig(cfg, savePath); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Config saved to %s\n", savePath)
	},
}

func sanitizeInput(input string) string {
	return string([]byte(input)[:len(input)-1]) // strip newline
}

func saveConfig(cfg InitConfig, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	return encoder.Encode(cfg)
}

func init() {
	rootCmd.AddCommand(initCmd)
}