package main

import (
	"log"

	"github.com/HollyEllmo/go_file_storage/p2p"
)

func main() {

tcpOpts := p2p.TCPTransportOpts{
	ListenAddr: ":3000",
	HandshakeFunc: p2p.NOPHandshakeFunc, // Default handshake function
	Decoder: p2p.DefaultDecoder{}, // Using GOB for decoding messages
}

tr := p2p.NewTCPTransport(tcpOpts)

if err := tr.ListenAndAccept(); err != nil {
	log.Fatalf("Failed to start TCP transport: %v", err)
}

select {}
}