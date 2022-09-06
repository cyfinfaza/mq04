package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
)

func main() {
	// Create the new MQTT Server.
	server := mqtt.NewServer(nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883")

	// Add the listener to the server with default options (nil).
	err := server.AddListener(tcp, nil)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	log.Println("Starting MQTT Server")
	// Start the broker. Serve() is blocking - see examples folder
	// for usage ideas.
	go server.Serve()
	<-done
	log.Println("Server stopped")
}
