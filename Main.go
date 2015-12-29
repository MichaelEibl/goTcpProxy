package main

import (
	"log"
	"os"
	"os/signal"
	"github.com/michaeleibl/goTcpProxy/config"
	"github.com/michaeleibl/goTcpProxy/tcp"
	"flag"
)

var quit = make(chan os.Signal, 1)

func main() {
	defer func() {
		log.Print("Trying to recover in main")
		if r := recover(); r != nil {
			log.Println("Recovered in f in main", r)
		} else {
			log.Println("No need to recover in main - all good")
		}

	}()

	flag.Parse()
	signal.Notify(quit, os.Interrupt, os.Kill)
	log.Println("Proxy Server Starting")
	config.LoadConfig("settings.xml")

	go startListeningServers()

	select {
	case <-quit:
	// received a kill or interrupt on the channel
		log.Printf("Application quitting\n")
		return
	}

}

func startListeningServers() {
	defer func() {
		log.Print("Trying to recover in startListeningServers")
		if err := recover(); err != nil {
			log.Println("Recovered in startListeningServers", err)
		} else {
			log.Println("No need to recover in startListeningServers - all good")
		}

	}()
	for _, proxyServer := range config.ProxyData.ProxyserverItems {
		safeStart(proxyServer)
	}
}

func safeStart(proxyServer config.Proxyserver) {
	log.Printf("Starting Listening Server %s", proxyServer.ProxyName)
	tcp.StartTCPListener(proxyServer)
}
