package config
import (
	"os"
	"io/ioutil"
	"encoding/xml"
	"log"
)



type Proxyserver struct {
	ProxyName       string `xml:"name,attr"`
	Buffersize      string `xml:"buffersize,attr"`
	SourceItem      Source `xml:"source"`
	DestinationItem Destination `xml:"destination"`
}

type Source struct {
	Port              string `xml:"port,attr"`
	Quedconnections   string `xml:"quedconnections,attr"`
	Receivebuffersize string `xml:"receivebuffersize,attr"`
	Sendbuffersize    string `xml:"sendbuffersize="`
	Bindaddress       string `xml:"bindaddress,attr"`
}

type Destination struct {
	Port              string `xml:"port,attr"`
	Ipaddress         string `xml:"ipaddress,attr"`
	Receivebuffersize string `xml:"receivebuffersize,attr"`
	Sendbuffersize    string `xml:"sendbuffersize,attr`
}

type Proxy struct {
	XMLName          xml.Name                       `xml:"proxy"`
	ProxyserverItems []Proxyserver `xml:"proxyserver"`
}

var ProxyData Proxy

func LoadConfig(fileName string) {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	var bytePayload []byte
	bytePayload, errReadAll := ioutil.ReadAll(xmlFile)
	if errReadAll != nil {
		log.Printf("error: %v", errReadAll)
		return
	}
	err2 := xml.Unmarshal(bytePayload, &ProxyData)
	if err2 != nil {
		log.Printf("error: %v", err2)
		panic("Proxy definition could not be unmarshalled. Please fix the source data")
		return
	}
	log.Println("XML data loaded ok")
	printData()

}

func printData() {
	for _, proxyServer := range ProxyData.ProxyserverItems {
		log.Printf("Name : %s", proxyServer.ProxyName)
		log.Printf("Buffersize : %s\n", proxyServer.Buffersize)
		log.Printf("\tSource\n")
		log.Printf("\t\tPort : %s\n", proxyServer.SourceItem.Port)
		log.Printf("\t\tBindaddress : %s\n", proxyServer.SourceItem.Bindaddress)
		log.Printf("\t\tQuedconnections : %s\n", proxyServer.SourceItem.Quedconnections)
		log.Printf("\t\tReceivebuffersize : %s\n", proxyServer.SourceItem.Receivebuffersize)
		log.Printf("\t\tSendbuffersize : %s\n", proxyServer.SourceItem.Sendbuffersize)
		log.Printf("\tDestination\n")
		log.Printf("\t\tPort : %s\n", proxyServer.DestinationItem.Port)
		log.Printf("\t\tReceivebuffersize : %s\n", proxyServer.DestinationItem.Receivebuffersize)
		log.Printf("\t\tSendbuffersize : %s\n", proxyServer.DestinationItem.Sendbuffersize)
		log.Printf("\t\tIpaddress : %s\n", proxyServer.DestinationItem.Ipaddress)
	}
	log.Println()
	log.Println()
	log.Println()
}