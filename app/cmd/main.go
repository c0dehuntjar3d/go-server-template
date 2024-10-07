package main

import (
	"app/pkg/initializer"
	"net/http"
)

func main() {
	initialize, _ := initializer.InitApplicaiton()

	server := initialize.Server

	server.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Warn("Info log")
	})

	server.Start()
}
