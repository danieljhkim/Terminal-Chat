package app

import "github.com/danieljhkim/chat-server/internal/protocol"

type Room struct {
	Name    string
	Members map[*Client]struct{}
}

// constructor
func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Members: make(map[*Client]struct{}),
	}
}

func (r *Room) Add(c *Client)    { r.Members[c] = struct{}{} }
func (r *Room) Remove(c *Client) { delete(r.Members, c) }

// Broadcast sends msg to every member, optional ‘skip’ (e.g. sender)
func (r *Room) Broadcast(msg protocol.WireMessage, skip *Client) {
	for m := range r.Members {
		if m == skip {
			continue
		}
		m.Send(msg)
	}
}
