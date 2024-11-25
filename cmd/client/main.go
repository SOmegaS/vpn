package main

import (
	"vpn/internal/client"
)

func main() {
	// Create new client
	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	// Init client
	err = c.Init()
	if err != nil {
		panic(err)
	}

	// Connect to another user
	c.Connect()
	//c.Listen()

	// Serve client
	err = c.Serve()
	if err != nil {
		panic(err)
	}
}
