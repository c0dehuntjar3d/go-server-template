package main

import (
	"fmt"
	"go-server/config"
	"go-server/pkg/initializer"
)

// http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
// 	server.Logger.Info("Info log", "test log")
// 	server.Logger.Debug("Debug log", "test log")
// 	server.Logger.Error("Error log", "test log")
// 	server.Logger.Warn("Warn log", "test log")
// })

func DefaultConfigStartup() {
	init := initializer.WithDefault()
	server := init.Server

	go server.Start()

	// waitForSignals(init.Logger, server)
	// shutdown(server, init.Logger)
}

func CustomConfigStartup() {

	cfg, err := config.New("config.yaml")
	if err != nil {
		panic(err)
	}

	init, err := initializer.New(cfg)
	if err != nil {
		fmt.Println(err)
	}

	server := init.Server
	go server.Start()

	// waitForSignals(init.Logger, server)
	// shutdown(server, init.Logger)
}
