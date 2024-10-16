package main

import (
	"app/pkg/initializer"
)

func main() {
	initialize, _ := initializer.InitApplicaiton()
	server := initialize.Server
	server.Start()
}
