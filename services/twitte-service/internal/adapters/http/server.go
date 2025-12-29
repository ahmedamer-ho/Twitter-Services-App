package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server represents HTTP delivery adapter
type Server struct {
	server *http.Server
}

// NewServer constructs HTTP server
func NewServer(handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         ":8082",
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// Run starts the server
func (s *Server) Run() error {
	log.Println("HTTP server listening on", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
