package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":8080"
	defaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	internalServer *http.Server
	channelErr     chan error
	// Timer for context cancellation
	shutdownTimeout time.Duration
}

func (s *Server) FullShutdownTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	log.Println("Shutting down server...")
	if err := s.internalServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Fatalf("Server shutdown successfully")
}

func (s *Server) GracefulShutdown() {
	osInterruptChan := make(chan os.Signal, 1)
	signal.Notify(osInterruptChan, os.Interrupt, os.Kill)
	select {
	case <-osInterruptChan:
		log.Fatal("Interruption by user or system\n")
	case err := <-s.channelErr:
		log.Fatalf("Server threw an error %v", err)
		s.FullShutdownTimeout()
	}
}

func (s *Server) Start() {
	go func() {
		s.channelErr <- s.internalServer.ListenAndServe()
		close(s.channelErr)
	}()
}

func NewServer(handler http.Handler, options ...Option) *Server {
	server := &Server{
		internalServer: &http.Server{
			Handler:      handler,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         defaultAddr,
		},
		channelErr:      make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}
	for _, option := range options {
		option(server)
	}

	server.Start()

	return server
}
