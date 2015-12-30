package tcp

const (
	SourceToDestination = iota
	DestinationToSource
)

type PacketFilter interface {
	Filter(data []byte, direction int)
}
