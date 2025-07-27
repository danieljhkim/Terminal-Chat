package protocol

// Message kinds.
const (
	TypeJoin      = "join"
	TypeLeave     = "leave"
	TypeRoomMsg   = "room_msg"
	TypeDM        = "dm"
	TypeListRooms = "list_rooms" // request
	TypeRoomsList = "rooms_list" // response
	TypeRoomsName = "rooms_name" // request
	TypeEcho      = "echo"
	TypeError     = "error" // error response
)

// WireMessage represents every payload on the wire.
type WireMessage struct {
	Type     string   `json:"type"`               // required
	Room     string   `json:"room,omitempty"`     // join / room_msg
	Target   string   `json:"target,omitempty"`   // dm
	Body     string   `json:"body,omitempty"`     // text
	Username string   `json:"username,omitempty"` // sender
	Rooms    []string `json:"rooms,omitempty"`    // list response
	Message  string   `json:"message,omitempty"`  // error message
}
