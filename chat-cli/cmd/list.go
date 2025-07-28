package cmd

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/danieljhkim/chat-cli/internal/config"
	"github.com/danieljhkim/chat-cli/internal/protocol"
	"github.com/spf13/cobra"
)

// roomsListCmd represents the list command
var roomsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available chat rooms",
	Long:    "Display all available chat rooms on the server.",
	Example: "chat-cli rooms list",
	RunE:    runListCommand,
}

// runListCommand handles the main logic for listing rooms
func runListCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	rooms, err := fetchRoomsList(cfg)
	if err != nil {
		return err
	}

	displayRooms(rooms)
	return nil
}

// fetchRoomsList connects to server and retrieves rooms list
func fetchRoomsList(cfg *config.Config) ([]string, error) {
	conn, err := establishConnection(cfg.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	if err := sendListRequest(enc, cfg.Username); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	rooms, err := receiveRoomsResponse(dec)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	return rooms, nil
}

// establishConnection creates a TCP connection to the server
func establishConnection(serverAddress string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", serverAddress, err)
	}
	return conn, nil
}

// sendListRequest sends the list rooms request to the server
func sendListRequest(enc *json.Encoder, username string) error {
	req := protocol.WireMessage{
		Type:     protocol.TypeListRooms,
		Username: username,
	}
	return enc.Encode(req)
}

// receiveRoomsResponse receives and validates the server response
func receiveRoomsResponse(dec *json.Decoder) ([]string, error) {
	var resp protocol.WireMessage
	if err := dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := validateResponse(&resp); err != nil {
		return nil, err
	}

	return resp.Rooms, nil
}

// validateResponse validates the server response
func validateResponse(resp *protocol.WireMessage) error {
	switch resp.Type {
	case protocol.TypeRoomsList:
		return nil
	case protocol.TypeError:
		return fmt.Errorf("server error: %s", resp.Message)
	default:
		return fmt.Errorf("unexpected response type: %s", resp.Type)
	}
}

// displayRooms formats and displays the rooms list
func displayRooms(rooms []string) {
	if len(rooms) == 0 {
		fmt.Println("No rooms available.")
		return
	}

	fmt.Println("Available rooms:")
	for i, room := range rooms {
		fmt.Printf("  %d. %s\n", i+1, room)
	}
	fmt.Printf("\nTotal: %d room(s)\n", len(rooms))
}

func init() {
	// Register the list command with the parent rooms command
	// This would be called from the parent package
	roomsCmd.AddCommand(roomsListCmd)
}
