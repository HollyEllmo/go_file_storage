package main

import (
	"log"

	"github.com/HollyEllmo/go_file_storage/p2p"
)

func main() {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    "localhost:3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: 	 p2p.DefaultDecoder{},
		// TODO: onPeer func
	}

	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot: "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport: tcpTransport,
		BootstrapNodes: []string{
			"localhost:4000",
		},
	}

	s := NewFileServer(fileServerOpts)

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	s.Stop()
	// }()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	
}