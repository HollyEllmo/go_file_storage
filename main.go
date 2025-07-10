package main

import (
	"fmt"
	"log"

	"github.com/HollyEllmo/go_file_storage/p2p"
)

func onPeer(peer p2p.Peer) error {
	peer.Close()
	// fmt.Println("doing some logic with the peer outside of TCPTransport")
	return nil
}

func main() {

tcpOpts := p2p.TCPTransportOpts{
	ListenAddr: ":3000",
	HandshakeFunc: p2p.NOPHandshakeFunc, // Default handshake function
	Decoder: p2p.DefaultDecoder{}, // Using GOB for decoding messages
	OnPeer: onPeer,
}

tr := p2p.NewTCPTransport(tcpOpts)

go func() {
	for {
		msg := <-tr.Consume()
		fmt.Printf("Received message from %+v\n", msg)
	}
}()

if err := tr.ListenAndAccept(); err != nil {
	log.Fatalf("Failed to start TCP transport: %v", err)
}

select {}
}