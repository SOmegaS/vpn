package main

import (
	"log"
	"os"
	"vpn/internal/client"
)

func tmplog() {
	file, err := os.OpenFile("vpn.log", os.O_CREATE, 0666) // os.O_APPEND
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
}

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
