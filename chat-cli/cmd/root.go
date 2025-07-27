/*
Copyright © 2025 Daniel Kim
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "chat-cli",
	Short: "A real-time command-line chat application",
	Long: `Chat CLI is a command-line interface for real-time chat communication.

Connect to chat servers, join rooms, and participate in conversations
directly from your terminal. Features include:
- Join, create and leave chat rooms
- Send and receive messages in real-time
- List available rooms

Example usage:
  chat-cli init,
  chat-cli rooms list
  chat-cli rooms join general`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {

	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}
