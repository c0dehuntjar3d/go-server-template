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

	"github.com/google/uuid"
)

type Server struct {
	Name            string
	Version         string
	Mux             *http.ServeMux
	Server          *http.Server
	Logger          logger.Interface
	notify          chan error
	shutdownTimeout time.Duration
}

func loggingMiddleware(next http.Handler, logger logger.Interface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		UUID := uuid.New().String()

		logger.Info(
			fmt.Sprintf(
				"Request: [%s] -> Path: [%s] | UUID: %s",
				r.Method,
				r.RequestURI,
				UUID,
			),
		)

		next.ServeHTTP(w, r)

		logger.Debug(
			fmt.Sprintf(
				"Request Completed: [%s] -> Path: [%s] in [%v] | UUID: %s",
				r.Method,
				r.RequestURI,
				time.Since(start),
				UUID,
			),
		)
	})
}

func New(cfg *config.HTTP, cfgApp *config.App, logger logger.Interface) *Server {
	mux := http.NewServeMux()
	server := &Server{
		Name:    cfgApp.Name,
		Version: cfgApp.Version,
		Logger:  logger,
		Mux:     mux,
		Server: &http.Server{
			Handler:      loggingMiddleware(mux, logger),
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
		s.Logger.Info(fmt.Sprintf("%s [%s] was started on port %s", s.Name, s.Version, s.Server.Addr))

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
