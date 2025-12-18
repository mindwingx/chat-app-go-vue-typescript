package http

import (
	"chat-app/domain/port"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Server struct {
	logger  port.ILogger
	server  *http.Server
	err     error
	SigChan chan os.Signal
}

func NewServer(port string, handler http.Handler, logger port.ILogger) *Server {
	corsHandler := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:4173", "http://localhost:5173", "http://localhost:9090"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      corsHandler.Handler(handler),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		SigChan: make(chan os.Signal, 1),
	}
}

func (s *Server) Start() {
	go func() {
		if s.err = s.server.ListenAndServe(); s.err != nil && !errors.Is(s.err, http.ErrServerClosed) {
			close(s.SigChan)
			s.logger.Error("server.start", zap.Error(s.err))
			return
		}
	}()

	log.Printf("[server] starting on port %s\n", s.server.Addr)
	return
}

func (s *Server) Stop(ctx context.Context) {
	if s.err != nil {
		log.Println("[server] failure!")
		return
	}

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("server.stop", zap.Error(err))
	}

	s.logger.Warn("service.stopped")
}
