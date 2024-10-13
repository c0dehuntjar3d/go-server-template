package main

import (
	rest "app/internal/controller/rest/user"
	repository "app/internal/repository/postgres/user"
	"app/pkg/initializer"
	"app/service/user"
	"net/http"
)

func main() {
	initialize, _ := initializer.InitApplicaiton()
	server := initialize.Server

	server.Mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Warn("Info log")
	})

	initUserCase(initialize)
	server.Start()
}

func initUserCase(initialize *initializer.Initializer) {
	repo, err := repository.NewUserRepository(initialize.DB, initialize.Logger)
	if err != nil {
		//todo
	}

	service, err := user.NewUserService(repo, initialize.Logger)
	if err != nil {
		//todo
	}
	rest.NewUserHandler(service, initialize.Logger, initialize.Server.Mux)
}
