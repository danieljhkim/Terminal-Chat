package protocol

import (
	"fmt"
	"time"
)

// Message types for different chat operations
const (
	// Room management
	TypeJoin    = "join"
	TypeLeave   = "leave"
	TypeRoomMsg = "room_msg"
	TypeAction  = "action" // /me style messages

	// Direct messaging
	TypeDM = "dm"
	TypeListDM = "list_dm" // request
	TypeDMList = "dm_list" // response
	TypeSendDM = "send_dm"

	// Room listing
	TypeListRooms = "list_rooms" // request
	TypeRoomsList = "rooms_list" // response
	TypeRoomsName = "rooms_name" // request

	// User management
	TypeListUsers  = "list_users"  // request
	TypeUserList   = "user_list"   // response
	TypeUserJoined = "user_joined" // notification
	TypeUserLeft   = "user_left"   // notification
	TypeUserCount  = "user_count"  // user count update

	// System messages
	TypeEcho    = "echo"
	TypePing    = "ping"    // heartbeat
	TypePong    = "pong"    // heartbeat response
	TypeStatus  = "status"  // server status
	TypeError   = "error"   // error response
	TypeInfo    = "info"    // informational message
	TypeWarning = "warning" // warning message

	// Server management
	TypeServerInfo = "server_info" // server information
	TypeStats      = "stats"       // server statistics
)

// WireMessage represents every payload on the wire
type WireMessage struct {
	// Core fields
	Type      string    `json:"type"`                // required - message type
	Timestamp time.Time `json:"timestamp,omitempty"` // message timestamp

	// Room and user identification
	Room     string `json:"room,omitempty"`     // room name for join/room_msg
	Username string `json:"username,omitempty"` // sender username
	Target   string `json:"target,omitempty"`   // target user for DM

	// Message content
	Body    string `json:"body,omitempty"`    // message text content
	Message string `json:"message,omitempty"` // system/error message

	// List responses
	Rooms []string `json:"rooms,omitempty"` // room list response
	Users []string `json:"users,omitempty"` // user list response
	DMs []DM `json:"DM,omitempty"` // user list response

	// Statistics and metadata
	UserCount    int               `json:"user_count,omitempty"`    // number of users in room
	MessageCount int               `json:"message_count,omitempty"` // total messages
	RoomCount    int               `json:"room_count,omitempty"`    // total rooms
	ServerUptime string            `json:"server_uptime,omitempty"` // server uptime
	Metadata     map[string]string `json:"metadata,omitempty"`      // additional data
}

type DM struct {
	Sender string `json:"sender"` // sender username
	Recipient string `json:"recipient"` // recipient username
	Body     string `json:"body"`     // message text content
	TimeStamp time.Time `json:"timestamp"` // message timestamp
}

// NewMessage creates a new WireMessage with timestamp
func NewMessage(msgType string) *WireMessage {
	return &WireMessage{
		Type:      msgType,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// NewRoomMessage creates a room message
func NewRoomMessage(room, username, body string) *WireMessage {
	msg := NewMessage(TypeRoomMsg)
	msg.Room = room
	msg.Username = username
	msg.Body = body
	return msg
}

// NewActionMessage creates an action message (/me style)
func NewActionMessage(room, username, action string) *WireMessage {
	msg := NewMessage(TypeAction)
	msg.Room = room
	msg.Username = username
	msg.Body = action
	return msg
}

// NewDirectMessage creates a direct message
func NewDirectMessage(sender, target, body string) *WireMessage {
	msg := NewMessage(TypeDM)
	msg.Username = sender
	msg.Target = target
	msg.Body = body
	return msg
}

// NewJoinMessage creates a join room message
func NewJoinMessage(room, username string) *WireMessage {
	msg := NewMessage(TypeJoin)
	msg.Room = room
	msg.Username = username
	return msg
}

// NewLeaveMessage creates a leave room message
func NewLeaveMessage(room, username string) *WireMessage {
	msg := NewMessage(TypeLeave)
	msg.Room = room
	msg.Username = username
	return msg
}

// NewUserJoinedNotification creates a user joined notification
func NewUserJoinedNotification(room, username string, userCount int) *WireMessage {
	msg := NewMessage(TypeUserJoined)
	msg.Room = room
	msg.Username = username
	msg.UserCount = userCount
	return msg
}

// NewUserLeftNotification creates a user left notification
func NewUserLeftNotification(room, username string, userCount int) *WireMessage {
	msg := NewMessage(TypeUserLeft)
	msg.Room = room
	msg.Username = username
	msg.UserCount = userCount
	return msg
}

// NewUserListResponse creates a user list response
func NewUserListResponse(room string, users []string) *WireMessage {
	msg := NewMessage(TypeUserList)
	msg.Room = room
	msg.Users = users
	msg.UserCount = len(users)
	return msg
}

// NewRoomsListResponse creates a rooms list response
func NewRoomsListResponse(rooms []string) *WireMessage {
	msg := NewMessage(TypeRoomsList)
	msg.Rooms = rooms
	msg.RoomCount = len(rooms)
	return msg
}

// NewErrorMessage creates an error message
func NewErrorMessage(errorMsg string) *WireMessage {
	msg := NewMessage(TypeError)
	msg.Message = errorMsg
	return msg
}

// NewInfoMessage creates an informational message
func NewInfoMessage(info string) *WireMessage {
	msg := NewMessage(TypeInfo)
	msg.Message = info
	return msg
}

// NewWarningMessage creates a warning message
func NewWarningMessage(warning string) *WireMessage {
	msg := NewMessage(TypeWarning)
	msg.Message = warning
	return msg
}

// NewServerStatsMessage creates a server statistics message
func NewServerStatsMessage(userCount, roomCount, messageCount int, uptime string) *WireMessage {
	msg := NewMessage(TypeStats)
	msg.UserCount = userCount
	msg.RoomCount = roomCount
	msg.MessageCount = messageCount
	msg.ServerUptime = uptime
	return msg
}

// NewPingMessage creates a ping message for heartbeat
func NewPingMessage() *WireMessage {
	return NewMessage(TypePing)
}

// NewPongMessage creates a pong response for heartbeat
func NewPongMessage() *WireMessage {
	return NewMessage(TypePong)
}

// IsSystemMessage returns true if the message is a system message
func (m *WireMessage) IsSystemMessage() bool {
	systemTypes := []string{
		TypeUserJoined, TypeUserLeft, TypeError, TypeInfo,
		TypeWarning, TypePing, TypePong, TypeStatus,
	}

	for _, sysType := range systemTypes {
		if m.Type == sysType {
			return true
		}
	}
	return false
}

// IsUserMessage returns true if the message is from a user
func (m *WireMessage) IsUserMessage() bool {
	userTypes := []string{TypeRoomMsg, TypeAction, TypeDM}

	for _, userType := range userTypes {
		if m.Type == userType {
			return true
		}
	}
	return false
}

// IsRequestMessage returns true if the message is a request
func (m *WireMessage) IsRequestMessage() bool {
	requestTypes := []string{
		TypeJoin, TypeLeave, TypeListRooms, TypeListUsers,
		TypeRoomsName, TypePing,
	}

	for _, reqType := range requestTypes {
		if m.Type == reqType {
			return true
		}
	}
	return false
}

// IsResponseMessage returns true if the message is a response
func (m *WireMessage) IsResponseMessage() bool {
	responseTypes := []string{
		TypeRoomsList, TypeUserList, TypePong, TypeError,
		TypeInfo, TypeStats,
	}

	for _, respType := range responseTypes {
		if m.Type == respType {
			return true
		}
	}
	return false
}

// SetMetadata adds metadata to the message
func (m *WireMessage) SetMetadata(key, value string) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
	m.Metadata[key] = value
}

// GetMetadata retrieves metadata from the message
func (m *WireMessage) GetMetadata(key string) (string, bool) {
	if m.Metadata == nil {
		return "", false
	}
	value, exists := m.Metadata[key]
	return value, exists
}

// Validate checks if the message has required fields
func (m *WireMessage) Validate() error {
	if m.Type == "" {
		return fmt.Errorf("message type is required")
	}

	switch m.Type {
	case TypeJoin, TypeLeave:
		if m.Room == "" || m.Username == "" {
			return fmt.Errorf("room and username are required for %s", m.Type)
		}
	case TypeRoomMsg, TypeAction:
		if m.Room == "" || m.Username == "" || m.Body == "" {
			return fmt.Errorf("room, username, and body are required for %s", m.Type)
		}
	case TypeDM:
		if m.Username == "" || m.Target == "" || m.Body == "" {
			return fmt.Errorf("username, target, and body are required for DM")
		}
	case TypeListUsers:
		if m.Room == "" {
			return fmt.Errorf("room is required for list users request")
		}
	}

	return nil
}
