package main

import (
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
	err = os.Chmod("vpn.log", 0777)
	if err != nil {
		panic(err)
	}
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
	log.Println("INFO: Mode connect")
	err = c.Connect()
	if err != nil {
		log.Panic(err)
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
