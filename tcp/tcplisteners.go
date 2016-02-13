package tcp

import (
	"log"
	"net"
	"strconv"
	"github.com/michaeleibl/tcpproxy/config"
	"fmt"
	"time"
)

var closeComs = make(chan net.Conn)

type TCPListenerMetaData struct {
	proxyServer config.Proxyserver
	filter      PacketFilter
}

func Init() {
	go closeSockets();
}

func closeSockets() {
	for {
		select {
		case coms := <-closeComs:
			if coms != nil {
				coms.Close()
			}

		}
	}
}

func StartTCPListener(proxyServer config.Proxyserver, filter PacketFilter) {
	tcpListener := createTcpListener(proxyServer, filter)
	go runTCPListener(tcpListener)
}

func createTcpListener(proxyServer config.Proxyserver, filter PacketFilter) *TCPListenerMetaData {
	tcpTempListner := &TCPListenerMetaData{proxyServer,
		filter, }

	return tcpTempListner
}

func runTCPListener(tcpListenerMetaData *TCPListenerMetaData) {
	defer func() {
		log.Print(">>>>>>>>>>>>>>>>>>Recover function call to check panic in runTCPListener")
		if r := recover(); r != nil {
			log.Println("Recovered in runTCPListener", r)
		} else {
			log.Println("Everything was fine in runTCPListener")
		}

	}()
	log.Printf("Opening Listener for proxy : %s\n", tcpListenerMetaData.proxyServer.ProxyName)
	service := fmt.Sprintf("%s:%s", tcpListenerMetaData.proxyServer.SourceItem.Bindaddress, tcpListenerMetaData.proxyServer.SourceItem.Port)
	log.Printf("Service : %s", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Printf("Error in Resolve TCP [%s]", err.Error())
		return
	}
	log.Printf("BEFORE - Binding listener IP[%s] and Port[%d]", tcpAddr.IP.String(), tcpAddr.Port)
	tcpListener, errListen := net.ListenTCP("tcp", tcpAddr)
	if errListen != nil {
		//log.Printf("Could not bind listener error [%s]", errListen.Error())
		panic(fmt.Sprintf("Could not bind listener error [%s]", errListen.Error()))
		//return
	}

	log.Printf("AFTER - Binding listener IP[%s] and Port[%d]", tcpAddr.IP.String(), tcpAddr.Port)

	for {
		sourceConn, errAccept := tcpListener.Accept()
		if errAccept != nil {
			log.Fatalln(errAccept)
		}
		log.Printf("Accepted connection : %+v\n", sourceConn.RemoteAddr().String())
		destinationConn := createDestination(tcpListenerMetaData)

		readTime := time.Time{}
		readTime.Add(time.Millisecond * 1000)
		errSetDeadlineSource := sourceConn.SetDeadline(readTime)
		if errSetDeadlineSource != nil {
			log.Fatalf("Cannot set deadline on socket %s ", errSetDeadlineSource)
		}
		errSetDeadlineDestination := destinationConn.SetDeadline(readTime)
		if errSetDeadlineDestination != nil {
			log.Fatalf("Cannot set deadline on socket %s ", errSetDeadlineDestination)
		}

		go sendFromDestinationToSource(sourceConn, destinationConn, tcpListenerMetaData)
		go sendFromSourceToDestination(sourceConn, destinationConn, tcpListenerMetaData)
	}

}

func createDestination(tcpListenerMetaData *TCPListenerMetaData) net.Conn {
	service := fmt.Sprintf("%s:%s", tcpListenerMetaData.proxyServer.DestinationItem.Ipaddress, tcpListenerMetaData.proxyServer.DestinationItem.Port)
	conn, err := net.DialTimeout("tcp", service, time.Millisecond * 1500)
	if err != nil {
		log.Fatal("Cannot connect to distination host ", err)
	}
	log.Printf("Connected to destination host %s on port %s", tcpListenerMetaData.proxyServer.DestinationItem.Ipaddress, tcpListenerMetaData.proxyServer.DestinationItem.Port)
	return conn
}

func closeConnections(sourceConn, destinationConn net.Conn) {
	closeComs <- sourceConn
	closeComs <- destinationConn
}

func sendFromSourceToDestination(sourceConn, destinationConn net.Conn, tcpListenerMetaData *TCPListenerMetaData) {
	defer closeConnections(sourceConn, destinationConn)
	value, conversionError := strconv.Atoi(tcpListenerMetaData.proxyServer.SourceItem.Receivebuffersize)
	if conversionError != nil {
		log.Fatalf("Buffer receive is wrong %s ", tcpListenerMetaData.proxyServer.SourceItem.Receivebuffersize)
	}
	var buf = make([]byte, value)


	//time.Sleep(time.Millisecond * 2500)

	for {
		// read the bytes into the buffer
		readLen, errRead := sourceConn.Read(buf[0:])
		if errRead != nil {
			log.Printf("Source Read error %s Data Read %d", errRead, readLen)
			if errRead.Error() == "EOF" || readLen == 0 {
				return
			}
		}
		// call the filter
		if tcpListenerMetaData != nil && tcpListenerMetaData.filter != nil {
			tcpListenerMetaData.filter.Filter(buf[0:readLen], SourceToDestination)
		}
		writeLen, errWrite := destinationConn.Write(buf[0:readLen])
		if errWrite != nil {
			log.Printf("Destination write error %s ", errWrite)
			return
		}
		// Add this for debug
		if *config.DebugFlag {
			log.Printf("Bytes read write %d:%d direction source->destination", readLen, writeLen)
		}
	}
}

func sendFromDestinationToSource(sourceConn, destinationConn net.Conn, tcpListenerMetaData *TCPListenerMetaData) {
	defer closeConnections(sourceConn, destinationConn)
	value, conversionError := strconv.Atoi(tcpListenerMetaData.proxyServer.DestinationItem.Receivebuffersize)
	if conversionError != nil {
		log.Fatalf("Buffer receive is wrong %s ", tcpListenerMetaData.proxyServer.SourceItem.Receivebuffersize)
	}
	var buf = make([]byte, value)
	for {
		// read the bytes into the buffer
		readLen, errRead := destinationConn.Read(buf[0:])
		if errRead != nil {
			log.Printf("Destination read error %s ", errRead)
			if errRead.Error() == "EOF" || readLen == 0 {
				return
			}
			return
		}
		// call the filter
		if tcpListenerMetaData != nil && tcpListenerMetaData.filter != nil {
			tcpListenerMetaData.filter.Filter(buf[0:readLen], DestinationToSource)
		}
		writeLen, errWrite := sourceConn.Write(buf[0:readLen])
		if errWrite != nil {
			log.Printf("Source write error %s", errWrite)
			return
		}
		if *config.DebugFlag {
			log.Printf("Bytes read write %d:%d direction destination->source", readLen, writeLen)
		}

	}
}

