package main

import (
	"fmt"
	"go-server/config"
	"go-server/pkg/httpserver"
	"go-server/pkg/initializer"
	"go-server/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.New("test_config.yaml")
	if err != nil {
		panic(err)
	}

	init, err := initializer.New(cfg)
	if err != nil {
		fmt.Println(err)
	}

	server := init.Server
	go server.Start()

	waitForSignals(init.Logger, server)
	shutdown(server, init.Logger)
}

func waitForSignals(log logger.Interface, httpServer *httpserver.Server) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}
	return err
}

func shutdown(httpServer *httpserver.Server, log logger.Interface) {
	err := httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
