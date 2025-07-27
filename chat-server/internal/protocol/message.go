package protocol

// Message kinds we understand.
const (
	TypeJoin    = "join"
	TypeLeave   = "leave"
	TypeRoomMsg = "room_msg"
	TypeDM      = "dm"
	TypeEcho    = "echo" // test / keep-alive
)

// WireMessage is the JSON payload on the wire.
type WireMessage struct {
	Type     string `json:"type"`               // required
	Room     string `json:"room,omitempty"`     // for join / room_msg
	Target   string `json:"target,omitempty"`   // for dm
	Body     string `json:"body,omitempty"`     // text
	Username string `json:"username,omitempty"` // sender’s handle (filled by client)
}
