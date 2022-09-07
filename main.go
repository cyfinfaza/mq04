package main

import (
	"flag"
	"log"
	"mqtt-server/utils"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/tg123/go-htpasswd"
)

func main() {
	log.Println("Initializing MQTT server")
	log.Println("Reading auth file")
	httpAuthFile := flag.String("auth", "./auth.htpasswd", "Path to htpasswd file for HTTP auth")
	needHelp := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *needHelp {
		flag.PrintDefaults()
		return
	}
	authDb, err := htpasswd.New(*httpAuthFile, htpasswd.DefaultSystems, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Configuring server")
	// Create the new MQTT Server.
	server := mqtt.NewServer(&mqtt.Options{
		BufferSize: 4096 * 1024,
	})

	listenerConfig := &listeners.Config{
		Auth: &utils.FileAuth{
			Checker: authDb,
		},
	}

	configuredListeners := []listeners.Listener{
		listeners.NewTCP("tcp0", ":1884"),
		listeners.NewWebsocket("ws0", ":8084"),
	}

	for _, l := range configuredListeners {
		log.Println("Adding listener", l.ID())
		server.AddListener(l, listenerConfig)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Start the broker. Serve() is blocking - see examples folder
	// for usage ideas.
	log.Println("Starting server")
	go server.Serve()
	log.Println("Serving")
	<-done
	log.Println("Server stopped")
}
