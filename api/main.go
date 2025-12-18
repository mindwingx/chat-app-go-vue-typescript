package main

import (
	"chat-app/delivery/http"
	"chat-app/di"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := di.ProvideLogger()
	defer logger.Stop()

	ws := di.ProvideWebSocketService(logger)
	{
		go ws.Broadcast(ctx)
		go ws.RetryFailedItems(ctx)
		go ws.ReleaseTypingUsers(ctx)
	}

	router := http.NewRouter(logger, ws)
	router.SetRoutes()

	server := http.NewServer("8080", http.RecoveryMiddleware(logger, router.Mux()), logger)
	server.Start()

	signal.Notify(server.SigChan, os.Interrupt, syscall.SIGTERM)
	<-server.SigChan

	logger.Warn("service.shutting-down")

	cancel()

	ws.Terminate()
	server.Stop(ctx)
}
