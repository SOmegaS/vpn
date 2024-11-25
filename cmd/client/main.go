package main

import (
	"fmt"
	"log"
	"os"
	"vpn/internal/client"
)

func main() {
	// Forward logs
	file, err := os.Create("vpn.log")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	log.SetOutput(file)

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
	fmt.Println("Choose mode (listen, connect): ")
	var mode string
	_, err = fmt.Scanln(&mode)
	if err != nil {
		panic(err)
	}
	if mode == "listen" {
		err = c.Listen()
		if err != nil {
			panic(err)
		}
	} else if mode == "connect" {
		err = c.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Serve client
	err = c.Serve()
	if err != nil {
		panic(err)
	}
}
