/*
Copyright © 2025 Daniel Kim
*/
package main

import (
	"fmt"

	"github.com/danieljhkim/chat-cli/cmd"
	"github.com/danieljhkim/chat-cli/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Configuration file not found or invalid.")
		cmd.PromptInitAndSave()
		return
	}
	fmt.Println("=== Chat-CLI Configuration ===")
	fmt.Printf("Server Address: %q\n", cfg.ServerAddress)
	fmt.Printf("Username: %q\n", cfg.Username)
	fmt.Println("==============================")
	fmt.Println()
	cmd.Execute()
}
