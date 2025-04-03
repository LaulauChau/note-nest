package http

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

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler, port int) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on %s", s.server.Addr)
		serverErrors <- s.server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received or an error occurs
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error starting server: %w", err)
	case <-shutdown:
		log.Println("Server is shutting down...")

		// Create a context with a timeout to wait for existing requests to complete
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Define a channel to signal when shutdown is done
		shutdownDone := make(chan struct{})

		// Shutdown the server gracefully
		go func() {
			defer close(shutdownDone)

			if err := s.server.Shutdown(ctx); err != nil {
				log.Printf("error during shutdown: %v", err)
			}
		}()

		// Wait for shutdown to complete or timeout
		select {
		case <-time.After(30 * time.Second):
			log.Println("server did not terminate gracefully, forcibly closing")
			if err := s.server.Close(); err != nil {
				log.Printf("error closing server: %v", err)
			}
		case <-shutdownDone:
			log.Println("server terminated gracefully")
		}
	}

	return nil
}
