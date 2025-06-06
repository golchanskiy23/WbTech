package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = "0.0.0.0:3333"
	defaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	internalServer  *http.Server
	channelErr      chan error
	shutdownTimeout time.Duration
}

func (s *Server) FullShutdownTimeout() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	log.Println("Shutting down server...")
	if err := s.internalServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown filed: %v", err)
	}
	log.Println("Server shutdown successfully")
	return nil
}

func (s *Server) GracefulShutdown() error {
	osInterruptChan := make(chan os.Signal, 1)
	signal.Notify(osInterruptChan, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-osInterruptChan:
		log.Printf("Server interrupted by system or user\n")
	case err := <-s.channelErr:
		log.Printf("Server threw an error %v\n", err)
	}
	close(osInterruptChan)
	if err := s.FullShutdownTimeout(); err != nil {
		return fmt.Errorf("graceful shutdown collapsed: %v", err)
	}
	return nil
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
