package server

import (
	"errors"
	"log"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	mux := http.NewServeMux()
	return &Server{
		mux: mux,
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("Starting server at %s", addr)
	err := http.ListenAndServe(addr, s.mux)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server error: %v", err)
		return err
	}
	log.Println("Server stopped gracefully.")
	return nil
}

func (s *Server) AddRoute(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}
