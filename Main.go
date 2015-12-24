package main
import (
	"log"
	"os"
	"os/signal"
	"fmt"
	"soffex.co.za/tcpproxy/config"
	"soffex.co.za/tcpproxy/tcp"
)

var quit = make(chan os.Signal, 1)

func main() {
	signal.Notify(quit, os.Interrupt, os.Kill)
	log.Println("Proxy Server Starting")
	config.LoadConfig("settings.xml")

	for _, proxyServer := range config.ProxyData.ProxyserverItems {
		log.Printf("Starting Listening Server %s", proxyServer.ProxyName)
		tcp.StartTCPListener(proxyServer)
	}

	select {
	case <-quit:
	// received a kill or interrupt on the channel
		fmt.Printf("Application quitting\n")
		return
	}

}
