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
	"time"

	"github.com/danieljhkim/chat-cli/internal/config"
	"github.com/danieljhkim/chat-cli/internal/protocol"
	"github.com/spf13/cobra"
)

var roomsJoinCmd = &cobra.Command{
    Use:     "join <room-name>",
    Short:   "Join a chat room",
    Long:    `Join a chat room and start participating in real-time conversations.`,
    Args:    cobra.MinimumNArgs(1),
    Example: "chat-cli rooms join general",
    RunE:    runJoinCommand,
}

// chatSession holds the state of the current chat session
type chatSession struct {
    roomName    string
    username    string
    conn        net.Conn
    enc         *json.Encoder
    dec         *json.Decoder
    userCount   int
    messageCount int
    startTime   time.Time
    showTimestamps bool
}

// runJoinCommand handles the main logic for joining a room
func runJoinCommand(cmd *cobra.Command, args []string) error {
    roomName := args[0]
    userName := ""
    if len(args) == 2 {
        userName = args[1]
    } else {
        cfg, err := config.Get()
        if err != nil {
            return fmt.Errorf("failed to load configuration: %w", err)
        }
        userName = cfg.Username
    }
    // Establish connection and join room
    conn, enc, dec, err := connectAndJoinRoom(roomName)
    if err != nil {
        return err
    }
    defer conn.Close()

    session := &chatSession{
        roomName:       roomName,
        username:       userName,
        conn:           conn,
        enc:            enc,
        dec:            dec,
        startTime:      time.Now(),
        showTimestamps: false,
    }
    printWelcome(session)
    return startAdvancedChatSession(session)
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

// printWelcome displays an enhanced welcome screen
func printWelcome(session *chatSession) {
    clearScreen()
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Printf("║                 🚀 Terminal Chat v1.0.0                  ║\n")
    fmt.Println("║══════════════════════════════════════════════════════════║")
    fmt.Printf("║ Room: %-51s║\n", session.roomName)
    fmt.Printf("║ User: %-51s║\n", session.username)
    fmt.Printf("║ Connected: %-46s║\n", session.startTime.Format("2006-01-02 15:04:05"))
    fmt.Println("║══════════════════════════════════════════════════════════║")
    fmt.Println("║                       COMMANDS                           ║")
    fmt.Println("║ /help      - Show this help                              ║")
    fmt.Println("║ /users     - List users in room                          ║")
    fmt.Println("║ /stats     - Show session statistics                     ║")
    fmt.Println("║ /time      - Toggle timestamps                           ║")
    fmt.Println("║ /clear     - Clear screen                                ║")
    fmt.Println("║ /quit      - Exit chat                                   ║")
    fmt.Println("║ Ctrl+C     - Exit gracefully                             ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()
    fmt.Printf("💬 Welcome to %s! Start typing to chat...\n", session.roomName)
    fmt.Println()
}

// startAdvancedChatSession manages the enhanced chat session
func startAdvancedChatSession(session *chatSession) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Setup graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Error channel for goroutine communication
    errChan := make(chan error, 2)

    // Start message listener
    go handleAdvancedIncomingMessages(ctx, session, errChan)

    // Start input handler
    go handleAdvancedUserInput(ctx, session, errChan)

    // Wait for either an error or interrupt signal
    select {
    case err := <-errChan:
        if err != nil {
            return fmt.Errorf("chat session error: %w", err)
        }
    case <-sigChan:
        fmt.Println("\n👋 Leaving room...")
        sendLeaveMessage(session)
    }

    cancel()
    return nil
}

// handleAdvancedIncomingMessages processes messages with enhanced formatting
func handleAdvancedIncomingMessages(ctx context.Context, session *chatSession, errChan chan<- error) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            var msg protocol.WireMessage
            if err := session.dec.Decode(&msg); err != nil {
                errChan <- fmt.Errorf("error reading message: %w", err)
                return
            }
            session.messageCount++
            switch msg.Type {
            case protocol.TypeRoomMsg:
                displayChatMessage(session, &msg)
            case protocol.TypeUserJoined:
                fmt.Printf("🟢 %s joined the room\n", msg.Username)
            case protocol.TypeUserLeft:
                fmt.Printf("🔴 %s left the room\n", msg.Username)
            case protocol.TypeUserList:
                displayUserList(msg.Users)
            case protocol.TypeError:
                fmt.Printf("❌ Server error: %s\n", msg.Message)
            default:
                fmt.Printf("❓ Unknown message type: %s\n", msg.Type)
            }
        }
    }
}

// displayChatMessage formats and displays a chat message
func displayChatMessage(session *chatSession, msg *protocol.WireMessage) {
    timestamp := ""
    if session.showTimestamps {
        timestamp = fmt.Sprintf("[%s] ", time.Now().Format("15:04:05"))
    }

    // Highlight own messages
    if msg.Username == session.username {
        // fmt.Printf("%s\033[36m[You]\033[0m: %s\n", timestamp, msg.Body)
    } else {
        fmt.Printf("%s\033[33m[%s]\033[0m: %s\n", timestamp, msg.Username, msg.Body)
    }
}

// displayUserList shows the list of users in the room
func displayUserList(userList []string) {
    fmt.Println("👥 Users in room:")
    for i, user := range userList {
        user = strings.TrimSpace(user)
        if user != "" {
            fmt.Printf("   %d. %s\n", i+1, user)
        }
    }
    fmt.Printf("Total: %d user(s)\n", len(userList))
}

// handleAdvancedUserInput processes user input with command support
func handleAdvancedUserInput(ctx context.Context, session *chatSession, errChan chan<- error) {
    scanner := bufio.NewScanner(os.Stdin)
    
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

            input := strings.TrimSpace(scanner.Text())
            if input == "" {
                continue
            }
            // Handle commands
            if strings.HasPrefix(input, "/") {
                if err := handleChatCommand(input, session); err != nil {
                    fmt.Printf("❌ Command error: %v\n", err)
                }
                continue
            }
            
            // Send regular message
            if err := sendChatMessage(input, session); err != nil {
                errChan <- fmt.Errorf("error sending message: %w", err)
                return
            }
        }
    }
}

// handleChatCommand processes chat commands
func handleChatCommand(input string, session *chatSession) error {
    parts := strings.Fields(input)
    if len(parts) == 0 {
        return fmt.Errorf("empty command")
    }

    command := strings.ToLower(parts[0])

    switch command {
    case "/help":
        printHelp()
    case "/quit", "/exit":
        fmt.Println("👋 Goodbye!")
        sendLeaveMessage(session)
        os.Exit(0)
    case "/clear":
        clearScreen()
        fmt.Printf("💬 Back in %s\n\n", session.roomName)
    case "/users":
        return requestUserList(session)
    case "/stats":
        printStats(session)
    case "/time":
        session.showTimestamps = !session.showTimestamps
        status := "disabled"
        if session.showTimestamps {
            status = "enabled"
        }
        fmt.Printf("🕒 Timestamps %s\n", status)
    case "/me":
        if len(parts) > 1 {
            action := strings.Join(parts[1:], " ")
            return sendActionMessage(action, session)
        }
        fmt.Println("Usage: /me <action>")
    default:
        fmt.Printf("❓ Unknown command: %s. Type /help for available commands.\n", command)
    }
    return nil
}

// printHelp displays the help menu
func printHelp() {
    fmt.Println("╔══════════════════════════════════════╗")
    fmt.Println("║              CHAT COMMANDS           ║")
    fmt.Println("╠══════════════════════════════════════╣")
    fmt.Println("║ /help      - Show this help          ║")
    fmt.Println("║ /users     - List users in room      ║")
    fmt.Println("║ /stats     - Show session stats      ║")
    fmt.Println("║ /time      - Toggle timestamps       ║")
    fmt.Println("║ /me <text> - Send action message     ║")
    fmt.Println("║ /clear     - Clear screen            ║")
    fmt.Println("║ /quit      - Exit chat               ║")
    fmt.Println("╚══════════════════════════════════════╝")
}

// printStats displays session statistics
func printStats(session *chatSession) {
    duration := time.Since(session.startTime)
    fmt.Println("📊 Session Statistics:")
    fmt.Printf("   Room: %s\n", session.roomName)
    fmt.Printf("   Connected for: %v\n", duration.Round(time.Second))
    fmt.Printf("   Messages received: %d\n", session.messageCount)
    fmt.Printf("   Started: %s\n", session.startTime.Format("2006-01-02 15:04:05"))
    requestUserList(session)
}

// clearScreen clears the terminal screen
func clearScreen() {
    fmt.Print("\033[2J\033[H")
}

// sendChatMessage sends a regular chat message
func sendChatMessage(text string, session *chatSession) error {
    msg := protocol.WireMessage{
        Type:     protocol.TypeRoomMsg,
        Room:     session.roomName,
        Body:     text,
        Username: session.username,
    }
    return session.enc.Encode(msg)
}

// sendActionMessage sends an action message (/me command)
func sendActionMessage(action string, session *chatSession) error {
    msg := protocol.WireMessage{
        Type:     protocol.TypeAction,
        Room:     session.roomName,
        Body:     action,
        Username: session.username,
    }
    return session.enc.Encode(msg)
}

// requestUserList requests the list of users in the room
func requestUserList(session *chatSession) error {
    msg := protocol.WireMessage{
        Type:     protocol.TypeListUsers,
        Room:     session.roomName,
        Username: session.username,
    }
    return session.enc.Encode(msg)
}

// sendLeaveMessage notifies the server that the user is leaving
func sendLeaveMessage(session *chatSession) {
    msg := protocol.WireMessage{
        Type:     protocol.TypeLeave,
        Room:     session.roomName,
        Username: session.username,
    }
    session.enc.Encode(msg) // Ignore error on shutdown
}

func init() {
    roomsCmd.AddCommand(roomsJoinCmd)
}