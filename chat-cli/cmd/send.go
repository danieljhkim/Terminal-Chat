/*
Copyright © 2025 Daniel Kim
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/danieljhkim/chat-cli/internal/config"
	"github.com/danieljhkim/chat-cli/internal/net"
	"github.com/danieljhkim/chat-cli/internal/protocol"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
    Use:   "send <username> <message>",
    Short: "Send a direct message to a specific user",
    Long: `Send a private direct message to another user on the server.
The message will only be visible to you and the recipient.

Example:
  chat-cli dm send alice Hello there! How are you today?`,
    Args: cobra.MinimumNArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        username := args[0]
        // Join all remaining args as the message content
        message := strings.Join(args[1:], " ")
        
        return sendDirectMessage(username, message)
    },
}

var cfg *config.Config

func init() {
	cf, err := config.Get()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
	} else {
		fmt.Println("Loaded configuration successfully.")
		cfg = cf
	}
    dmCmd.AddCommand(sendCmd)
}

// sendDirectMessage handles connecting to the server and sending a DM
func sendDirectMessage(targetUser, messageContent string) error {
    currentUser := cfg.Username    
    conn, err := net.Connect(cfg.ServerAddress)
    if err != nil {
        return fmt.Errorf("failed to connect to server: %w", err)
    }
    defer conn.Close()
    enc := json.NewEncoder(conn)
    
    dmMsg := protocol.WireMessage{
        Type:     protocol.TypeSendDM,
        Target:   targetUser,
        Body:     messageContent,
        Username: currentUser,
        Timestamp: time.Now(),
    }
    
    if err := enc.Encode(dmMsg); err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    return nil
}
