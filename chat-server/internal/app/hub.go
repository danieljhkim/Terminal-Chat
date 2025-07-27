package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danieljhkim/chat-server/internal/protocol"
)

type envelope struct {
	sender *Client
	msg    protocol.WireMessage
}

type Hub struct {
	Rooms      map[string]*Room
	Clients    map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	Inbound    chan envelope

	log *slog.Logger
}

func NewHub(log *slog.Logger) *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Inbound:    make(chan envelope, 1024),
		log:        log.With("component", "hub"),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = struct{}{}
			h.log.Info("client connected", "addr", c.conn.RemoteAddr())
		case c := <-h.Unregister:
			h.removeClient(c)
		case env := <-h.Inbound:
			h.dispatch(env)
		case <-ctx.Done():
			h.log.Info("hub shutting down")
			return
		}
	}
}

func (h *Hub) removeClient(c *Client) {
	for _, r := range h.Rooms {
		r.Remove(c)
	}
	delete(h.Clients, c)
	c.Close()
}

func (h *Hub) getOrCreateRoom(name string) *Room {
	if r, ok := h.Rooms[name]; ok {
		return r
	}
	r := NewRoom(name)
	h.Rooms[name] = r
	return r
}

/* -------------------------------------------------- *
 *                     Routing                        *
 * -------------------------------------------------- */

func (h *Hub) dispatch(env envelope) {
	msg := env.msg
	c := env.sender

	switch msg.Type {

	case protocol.TypeJoin:
		h.handleJoin(c, msg)

	case protocol.TypeRoomMsg:
		h.handleRoomMsg(c, msg)

	case protocol.TypeDM:
		h.handleDM(c, msg)

	case protocol.TypeEcho:
		c.Send(msg) // simple echo

	default:
		h.log.Warn("unknown msg type", "type", msg.Type)
	}
}

func (h *Hub) handleJoin(c *Client, msg protocol.WireMessage) {
	if msg.Room == "" {
		return
	}
	room := h.getOrCreateRoom(msg.Room)
	room.Add(c)

	notice := protocol.WireMessage{
		Type: protocol.TypeRoomMsg,
		Room: room.Name,
		Body: fmt.Sprintf("%s joined the room.", msg.Username),
	}
	room.Broadcast(notice, nil)
}

func (h *Hub) handleRoomMsg(c *Client, msg protocol.WireMessage) {
	if msg.Room == "" {
		return
	}
	if room, ok := h.Rooms[msg.Room]; ok {
		// let sender’s Username go through unchanged
		room.Broadcast(msg, nil)
	}
}

func (h *Hub) handleDM(_ *Client, msg protocol.WireMessage) {
	if msg.Target == "" {
		return
	}
	for cl := range h.Clients {
		if cl.Username == msg.Target {
			cl.Send(msg)
			return
		}
	}
}
