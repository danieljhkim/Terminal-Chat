/*
Copyright © 2025 Daniel Kim
*/
package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/danieljhkim/chat-cli/internal/config"
	"github.com/danieljhkim/chat-cli/internal/protocol"
	"github.com/spf13/cobra"
)

var roomsJoinCmd = &cobra.Command{
	Use:     "join <room-name>",
	Short:   "Join a chat room",
	Long:    `Join a chat room and start participating in real-time conversations.`,
	Args:    cobra.ExactArgs(1),
	Example: "chat-cli rooms join general",
	RunE:    runJoinCommand,
}

// runJoinCommand handles the main logic for joining a room
func runJoinCommand(cmd *cobra.Command, args []string) error {
	roomName := args[0]

	// Establish connection and join room
	conn, enc, dec, err := connectAndJoinRoom(roomName)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Print("==============================\n")
	fmt.Printf("Successfully joined room: %s\n", roomName)
	fmt.Println("Type your message and press Enter.")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Print("==============================\n")

	return startChatSession(roomName, enc, dec)
}

// connectAndJoinRoom establishes connection and sends join request
func connectAndJoinRoom(roomName string) (net.Conn, *json.Encoder, *json.Decoder, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	conn, err := net.Dial("tcp", cfg.ServerAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	// Send join request
	joinReq := protocol.WireMessage{
		Type:     protocol.TypeJoin,
		Room:     roomName,
		Username: cfg.Username,
	}
	if err := enc.Encode(joinReq); err != nil {
		conn.Close()
		return nil, nil, nil, fmt.Errorf("failed to send join request: %w", err)
	}

	// Read join response
	var resp protocol.WireMessage
	if err := dec.Decode(&resp); err != nil {
		conn.Close()
		return nil, nil, nil, fmt.Errorf("failed to read join response: %w", err)
	}
	if resp.Type == protocol.TypeError {
		conn.Close()
		return nil, nil, nil, fmt.Errorf("server error: %s", resp.Message)
	}

	return conn, enc, dec, nil
}

// startChatSession manages the chat session with concurrent message handling
func startChatSession(roomName string, enc *json.Encoder, dec *json.Decoder) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Error channel for goroutine communication
	errChan := make(chan error, 2)

	// Start message listener
	go handleIncomingMessages(ctx, dec, errChan)

	// Start input handler
	go handleUserInput(ctx, roomName, enc, errChan)

	// Wait for either an error or interrupt signal
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("chat session error: %w", err)
		}
	case <-sigChan:
		fmt.Println("\nLeaving room...")
	}

	cancel()
	return nil
}

// handleIncomingMessages processes messages received from the server
func handleIncomingMessages(ctx context.Context, dec *json.Decoder, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg protocol.WireMessage
			if err := dec.Decode(&msg); err != nil {
				errChan <- fmt.Errorf("error reading message: %w", err)
				return
			}

			switch msg.Type {
			case protocol.TypeRoomMsg:
				fmt.Printf("[%s]: %s\n", msg.Username, msg.Body)
			case protocol.TypeError:
				fmt.Printf("Server error: %s\n", msg.Message)
			default:
				fmt.Printf("Received unknown message type: %s\n", msg.Type)
			}
		}
	}
}

// handleUserInput processes user input and sends messages to the server
func handleUserInput(ctx context.Context, roomName string, enc *json.Encoder, errChan chan<- error) {
	scanner := bufio.NewScanner(os.Stdin)
	cfg, err := config.Get()
	if err != nil {
		errChan <- fmt.Errorf("failed to load configuration: %w", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					errChan <- fmt.Errorf("error reading input: %w", err)
				}
				return
			}

			text := strings.TrimSpace(scanner.Text())
			if text == "" {
				continue
			}

			msg := protocol.WireMessage{
				Type:     protocol.TypeRoomMsg,
				Room:     roomName,
				Body:     text,
				Username: cfg.Username,
			}

			if err := enc.Encode(msg); err != nil {
				errChan <- fmt.Errorf("error sending message: %w", err)
				return
			}
		}
	}
}

func init() {
	roomsCmd.AddCommand(roomsJoinCmd)
}
