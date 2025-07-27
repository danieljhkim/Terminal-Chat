package netutil

import (
	"context"
	"log/slog"
	"net"

	"github.com/danieljhkim/chat-server/internal/app"
)

func StartTCP(ctx context.Context, addr string, hub *app.Hub, log *slog.Logger) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("tcp listening", "addr", addr)

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Warn("accept error", "err", err)
				continue
			}
		}

		client := app.NewClient(conn, hub, log)
		hub.Register <- client

		go client.ReadLoop()
		go client.WriteLoop()
	}
}
