package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/danieljhkim/chat-cli/internal/config"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure the chat CLI",
	Long:  "Initialize the chat CLI by setting up server address and username configuration.",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	return PromptInitAndSave(args...)
}

func PromptInitAndSave(args ...string) error {
	cfg, err := promptForConfig()
	if err != nil {
		return fmt.Errorf("failed to collect configuration: %w", err)
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to determine config path: %w", err)
	}

	if err := config.Save(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Configuration saved to %s\n", configPath)
	return nil
}

func promptForConfig() (*config.Config, error) {
	reader := bufio.NewReader(os.Stdin)
	serverAddr, err := promptInput(reader, "Enter server address (e.g. localhost:9000): ")
	if err != nil {
		return nil, fmt.Errorf("failed to read server address: %w", err)
	}
	username, err := promptInput(reader, "Enter username: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read username: %w", err)
	}
	return &config.Config{
		ServerAddress: serverAddr,
		Username:      username,
	}, nil
}

func promptInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return sanitizeInput(input), nil
}

func sanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

func init() {
	rootCmd.AddCommand(InitCmd)
}
