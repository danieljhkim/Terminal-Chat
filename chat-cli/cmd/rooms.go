/*
Copyright © 2025 Daniel Kim
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var roomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "Manage chat rooms",
	Long: `Manage chat rooms including listing available rooms, joining rooms, 
and leaving rooms. Use subcommands to perform specific room operations.

Available subcommands:
  list  - List all available chat rooms
  join  - Join a specific chat room`,
	Example: `  chat-cli rooms list
  chat-cli rooms join general`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(roomsCmd)
}
