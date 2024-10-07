package httpserver

import (
	"app/config"
	"app/pkg/logger"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		s.Logger.Info(fmt.Sprint("Server was started on port: ", s.Server.Addr))

		s.notify <- s.Server.ListenAndServe()
		close(s.notify)
	}()

	s.waitForSignals()
	s.shutdown()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) waitForSignals() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case <-interrupt:
		s.Logger.Info("Application server is stopping by interrupt..")
	case err = <-s.Notify():
		s.Logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}
	return err
}

func (s *Server) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.Server.Shutdown(ctx)
	if err != nil {
		s.Logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
	s.Logger.Info("Application server is stopped by interrupt..")
}
