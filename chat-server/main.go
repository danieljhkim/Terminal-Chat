package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/danieljhkim/chat-server/internal/app"
	"github.com/danieljhkim/chat-server/internal/config"
	"github.com/danieljhkim/chat-server/internal/logger"
	netutil "github.com/danieljhkim/chat-server/internal/net"
)

func main() {
	if err := config.Load("config"); err != nil {
		panic(err)
	}

	log := logger.New(config.Cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	hub := app.NewHub(log)
	go hub.Run(ctx)

	if err := netutil.StartTCP(ctx, config.Cfg.ListenAddress, hub, log); err != nil {
		log.Error("listener error", "err", err)
	}
}
