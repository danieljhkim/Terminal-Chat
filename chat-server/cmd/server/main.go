package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/danieljhkim/chat-server/internal/app"
	"github.com/danieljhkim/chat-server/internal/config"
	"github.com/danieljhkim/chat-server/internal/logger"
	netutil "github.com/danieljhkim/chat-server/internal/net"
)

// main is the entry-point for the chat server binary.
// Usage: `go run ./cmd/server` (or build + run the binary produced).
func main() {
	// 1) Load configuration (from config.yaml, $CHAT_SERVER_* env vars, or defaults).
	if err := config.Load("config"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 2) Initialise a structured logger at the configured level.
	log := logger.New(config.Cfg.LogLevel)
	log.Info("configuration loaded",
		"listen_address", config.Cfg.ListenAddress,
		"transport", config.Cfg.Transport,
		"max_message_bytes", config.Cfg.MaxMessageBytes,
	)

	// 3) Graceful-shutdown context (Ctrl-C → cancel).
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 4) Create and run the Hub (central router) in its own goroutine.
	hub := app.NewHub(log)
	go hub.Run(ctx)

	// 5) Start a network listener based on the configured transport.
	var listenErr error
	switch config.Cfg.Transport {
	case "tcp":
		listenErr = netutil.StartTCP(ctx, config.Cfg.ListenAddress, hub, log)
	// case "websocket":
	//     listenErr = netutil.StartWebSocket(ctx, config.Cfg.ListenAddress, hub, log)
	default:
		log.Error("unknown transport", "transport", config.Cfg.Transport)
		os.Exit(1)
	}

	if listenErr != nil {
		log.Error("server stopped", "err", listenErr)
	}
}
