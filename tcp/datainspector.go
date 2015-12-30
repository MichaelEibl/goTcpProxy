package tcp

import (
	"github.com/michaeleibl/tcpproxy/config"
	"log"
)

type DataInspection struct {

}

func NewDataInspection() *DataInspection {
	return &DataInspection{}
}

func (d DataInspection) Filter(data []byte, direction int) {
	//TODO inspect the data here
	if *config.DebugFlag {
		switch direction {
		case SourceToDestination:
			log.Printf("Data inspection S->D [%s]",  string(data))
			break;
		case DestinationToSource:
			log.Printf("Data inspection S<-D [%s]",  string(data))
			break;
		}
	}
}
