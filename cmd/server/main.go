package main

import (
	"fmt"

	"app/internal/app/adapter"
	"app/internal/app/adapter/user"
	"app/internal/pkg/initializer"
)

func main() {
	initializr, cfg := initializer.InitApplication()
	server := initializr.Server

	adapters := []adapter.Adapter{
		user.NewUserAdapter(),
	}

	for _, a := range adapters {
		err := a.Initialize(initializr.Logger, *cfg, initializr.Server, initializr.DB)
		if err != nil {
			initializr.Logger.Fatal(fmt.Sprintf("Usecase %s load error: %s", a.Name(), err.Error()))
		} else {
			initializr.Logger.Debug(fmt.Sprintf("Usecase %s load success", a.Name()))
		}
	}

	server.Start()
}
