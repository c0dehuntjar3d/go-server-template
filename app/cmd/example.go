package main

import (
	"app/config"
	"app/pkg/initializer"
	"fmt"
)

func DefaultConfigStartup() {
	init := initializer.WithDefault()
	server := init.Server

	go server.Start()

	// waitForSignals(init.Logger, server)
	// shutdown(server, init.Logger)
}

func CustomConfigStartup() {

	cfg, err := config.New("example_config.yaml")
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
