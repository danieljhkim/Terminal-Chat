package app

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"net"
	"time"

	"github.com/danieljhkim/chat-server/internal/config"
	"github.com/danieljhkim/chat-server/internal/protocol"
)

type Client struct {
	Username string
	conn     net.Conn
	hub      *Hub
	send     chan protocol.WireMessage
	log      *slog.Logger
}

func NewClient(conn net.Conn, hub *Hub, log *slog.Logger) *Client {
	return &Client{
		conn: conn,
		hub:  hub,
		send: make(chan protocol.WireMessage, 256),
		log:  log.With("addr", conn.RemoteAddr()),
	}
}

func (c *Client) Close() {
	_ = c.conn.Close()
	close(c.send)
}

func (c *Client) Send(msg protocol.WireMessage) {
	select {
	case c.send <- msg:
	default:
		c.log.Warn("send buffer full, dropping message")
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.hub.Unregister <- c
	}()

	reader := bufio.NewReaderSize(c.conn, config.Cfg.MaxMessageBytes)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}
		var msg protocol.WireMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			c.log.Warn("bad json", "err", err)
			continue
		}
		if c.Username == "" && msg.Username != "" {
			c.Username = msg.Username
		}
		c.hub.Inbound <- envelope{sender: c, msg: msg}
	}
}

func (c *Client) WriteLoop() {
	enc := json.NewEncoder(c.conn)
	for msg := range c.send {
		_ = c.conn.SetWriteDeadline(time.Now().Add(config.Cfg.WriteTimeout))
		if err := enc.Encode(msg); err != nil {
			c.log.Warn("write failed", "err", err)
			return
		}
	}
}
