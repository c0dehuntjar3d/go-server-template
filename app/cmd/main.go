package main

import (
	"app/pkg/initializer"
	"net/http"
)

func main() {
	initialize, _ := initializer.InitApplicaiton()

	server := initialize.Server

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Info("Info log")
	})

	server.Start()
}
