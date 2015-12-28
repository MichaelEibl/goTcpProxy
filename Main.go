package main

import (
	"log"
	"os"
	"os/signal"
	"fmt"
	"github.com/MichaelEibl/goTcpProxy/config"
	"github.com/MichaelEibl/goTcpProxy/tcp"
	"flag"
)

var quit = make(chan os.Signal, 1)

func main() {
	defer func() {
		log.Print("Trying to recover in main")
		if r := recover(); r != nil {
			fmt.Println("Recovered in f in main", r)
		}

	}()

	flag.Parse()
	signal.Notify(quit, os.Interrupt, os.Kill)
	log.Println("Proxy Server Starting")
	config.LoadConfig("settings.xml")

	for _, proxyServer := range config.ProxyData.ProxyserverItems {
		safeStart(proxyServer)
	}

	select {
	case <-quit:
	// received a kill or interrupt on the channel
		fmt.Printf("Application quitting\n")
		return
	}

}

func safeStart(proxyServer config.Proxyserver) {
	defer func() {
		log.Print("Trying to recover")
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}

	}()

	log.Printf("Starting Listening Server %s", proxyServer.ProxyName)
	tcp.StartTCPListener(proxyServer)
}
