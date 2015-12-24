package tcp
import (
	"log"
	"net"
	"fmt"
	"strconv"
	"soffex.co.za/tcpproxy/config"
)


type TCPListenerMetaData struct {
	proxyServer config.Proxyserver
}


func StartTCPListener(proxyServer config.Proxyserver) {
	tcpListener := createTcpListener(proxyServer)
	go runTCPListener(tcpListener)
}

func createTcpListener(proxyServer config.Proxyserver) *TCPListenerMetaData {
	tcpTempListner := &TCPListenerMetaData{proxyServer, }


	return tcpTempListner
}

func runTCPListener(tcpListenerMetaData *TCPListenerMetaData) {

	log.Printf("Opening Listener for proxy : %s/n", tcpListenerMetaData.proxyServer.ProxyName)
	service := fmt.Sprintf("%s:%s", tcpListenerMetaData.proxyServer.SourceItem.Bindaddress, tcpListenerMetaData.proxyServer.SourceItem.Port)
	log.Printf("Service : %s", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Binding listener IP[%s] and Port[%d]", tcpAddr.IP.String(), tcpAddr.Port)
	tcpListener, errListen := net.ListenTCP("tcp", tcpAddr)
	if errListen != nil {
		log.Fatalln(errListen)
	}
	for {
		sourceConn, errAccept := tcpListener.Accept()
		if errAccept != nil {
			log.Fatalln(errAccept)
		}
		fmt.Printf("Accepted connection : %+v\n", sourceConn.RemoteAddr().String())
		destinationConn := createDestination(tcpListenerMetaData)

		go sendFromSourceToDestination(sourceConn, destinationConn, tcpListenerMetaData)
		go sendFromDestinationToSource(sourceConn, destinationConn, tcpListenerMetaData)
	}

}

func createDestination(tcpListenerMetaData *TCPListenerMetaData) net.Conn {
	service := fmt.Sprintf("%s:%s", tcpListenerMetaData.proxyServer.DestinationItem.Ipaddress, tcpListenerMetaData.proxyServer.DestinationItem.Port)
	conn, err := net.Dial("tcp", service)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func closeConnections(sourceConn, destinationConn net.Conn) {
	sourceConn.Close()
	destinationConn.Close()
}

func sendFromSourceToDestination(sourceConn, destinationConn net.Conn, tcpListenerMetaData *TCPListenerMetaData) {
	defer closeConnections(sourceConn, destinationConn)
	value, conversionError := strconv.Atoi(tcpListenerMetaData.proxyServer.SourceItem.Receivebuffersize)
	if conversionError != nil {
		log.Fatalf("Buffer receive is wrong %s ", tcpListenerMetaData.proxyServer.SourceItem.Receivebuffersize)
	}
	var buf = make([]byte, value)
	for {
		// read the bytes into the buffer
		readLen, errRead := sourceConn.Read(buf[0:])
		if errRead != nil {
			log.Println(errRead)
			return
		}
		//writeLen, errWrite := destinationConn.Write(buf[0:readLen])
		_, errWrite := destinationConn.Write(buf[0:readLen])
		if errWrite != nil {
			log.Println(errWrite)
			return
		}
		// Add this for debug
		//log.Printf("Bytes read write %d:%d direction source->destination", readLen, writeLen)
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
			log.Println(errRead)
			return
		}
		//writeLen, errWrite := sourceConn.Write(buf[0:readLen])
		_, errWrite := sourceConn.Write(buf[0:readLen])
		if errWrite != nil {
			log.Println(errWrite)
			return
		}
		//log.Printf("Bytes read write %d:%d direction destination->source", readLen, writeLen)
	}
}

