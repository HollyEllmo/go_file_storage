package p2p

// Message holds any arbitrary data that is being sent over the each transport between two noodes in the network.
type Message struct {
	Payload []byte
}