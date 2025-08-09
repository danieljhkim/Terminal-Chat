package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/danieljhkim/chat-cli/internal/config"
	"github.com/danieljhkim/chat-cli/internal/net"
	"github.com/danieljhkim/chat-cli/internal/protocol"
	"github.com/spf13/cobra"
)

// dmlist represents the list command
var dmListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available chat rooms",
	Long:    "Display all available chat rooms on the server.",
	Example: "chat-cli rooms list",
	RunE:    runDMListCommand,
}

func runDMListCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	dms, err := fetchDMList(cfg)
	if err != nil {
		return err
	}

	displayDM(dms)
	return nil
}

func fetchDMList(cfg *config.Config) ([]protocol.DM, error) {
	conn, err := net.Connect(cfg.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	if err := sendDMListRequest(enc, cfg.Username); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	dms, err := receiveDMListResponse(dec)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	return dms, nil
}

func sendDMListRequest(enc *json.Encoder, username string) error {
	req := protocol.WireMessage{
		Type:     protocol.TypeListDM,
		Username: username,
	}
	return enc.Encode(req)
}

func receiveDMListResponse(dec *json.Decoder) ([]protocol.DM, error) {
	var resp protocol.WireMessage
	if err := dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := validateDMlistResponse(&resp); err != nil {
		return nil, err
	}

	return resp.DMs, nil
}

func validateDMlistResponse(resp *protocol.WireMessage) error {
	switch resp.Type {
	case protocol.TypeDMList:
		return nil
	case protocol.TypeError:
		return fmt.Errorf("server error: %s", resp.Message)
	default:
		return fmt.Errorf("unexpected response type: %s", resp.Type)
	}
}

func displayDM(dms []protocol.DM) {
	if len(dms) == 0 {
		fmt.Println("No dms available.")
		return
	}

	fmt.Printf("\nTotal: %d DM(s)\n", len(dms))
	fmt.Println("Available DM's:")
	for _, dm := range dms {
		fmt.Printf("  [%s] %s: %s\n", dm.TimeStamp, dm.Sender, dm.Body)
	}
}

func init() {
	// Register the list command with the parent dms command
	// This would be called from the parent package
	dmCmd.AddCommand(dmListCmd)
}
