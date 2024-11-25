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
	log.Println("INFO: Initialized log")

	// Create new client
	log.Println("INFO: Creating client")
	c, err := client.NewClient()
	if err != nil {
		log.Panic(err)
	}
	log.Println("INFO: Created client")

	log.Println("INFO: Initializing client")
	// Init client
	err = c.Init()
	if err != nil {
		log.Panic(err)
	}
	log.Println("INFO: Initialized client")

	// Connect to another user
	log.Println("INFO: Choosing mode")
	fmt.Println("Choose mode (listen, connect): ")
	var mode string
	_, err = fmt.Scanln(&mode)
	if err != nil {
		log.Panic(err)
	}
	if mode == "listen" {
		log.Println("INFO: Mode listen")
		err = c.Listen()
		if err != nil {
			log.Panic(err)
		}
	} else if mode == "connect" {
		log.Println("INFO: Mode connect")
		err = c.Connect()
		if err != nil {
			log.Panic(err)
		}
	} else {
		log.Panic("wrong mode")
	}
	log.Println("INFO: Connected host")

	// Serve client
	log.Println("INFO: Serving")
	err = c.Serve()
	if err != nil {
		log.Panic(err)
	}
	log.Println("INFO: Exiting")
}
