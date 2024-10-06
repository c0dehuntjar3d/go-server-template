package httpserver

import (
	"app/config"
	"app/pkg/logger"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Server          *http.Server
	Logger          logger.Interface
	notify          chan error
	shutdownTimeout time.Duration
}

func New(cfg *config.HTTP, logger logger.Interface) *Server {
	server := &Server{
		Logger: logger,
		Server: &http.Server{
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			Addr:         net.JoinHostPort("", cfg.Address),
		},
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.ShutdownTimeout,
	}

	return server
}

func (s *Server) Start() {
	s.Logger.Info(fmt.Sprint("Server was started on port: ", s.Server.Addr))

	s.notify <- s.Server.ListenAndServe()
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
