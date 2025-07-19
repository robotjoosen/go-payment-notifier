package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	port   int
	mux    *http.ServeMux
	server *http.Server
}

func New(port int, endpoints map[string]http.HandlerFunc) *Server {
	srv := &Server{
		port: port,
	}

	srv.mux = http.NewServeMux()
	for pattern, handler := range endpoints {
		srv.mux.HandleFunc(pattern, handler)
	}

	return srv
}

func (s *Server) AddEndpoint(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) Run() {
	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server starting", "address", s.server.Addr, "port", s.port)

	if err := s.server.ListenAndServe(); err != nil {
		slog.Error("server stopped with error", "err", err.Error())

		return
	}

	slog.Info("server stopped")
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return
	}
}
