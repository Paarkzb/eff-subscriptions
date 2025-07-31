package HTTPServer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port int, timeout time.Duration, handler http.Handler) *Server {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
	}

	return &Server{httpServer: httpServer}
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
