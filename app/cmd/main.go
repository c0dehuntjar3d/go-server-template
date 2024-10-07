package main

import (
	"app/pkg/httpserver"
	"app/pkg/initializer"
	"app/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize, _ := initializer.InitApplicaiton()

	server := initialize.Server

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Info("Info log")
	})

	go server.Start()

	waitForSignals(server, initialize.Logger)
	shutdown(server, initialize.Logger)
}

func waitForSignals(httpServer *httpserver.Server, log logger.Interface) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
		shutdown(httpServer, log)
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
