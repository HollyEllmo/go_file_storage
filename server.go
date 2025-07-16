package main

import (
	"fmt"
	"log"

	"github.com/HollyEllmo/go_file_storage/p2p"
)

type FileServerOpts struct {
	StorageRoot string
	PathTransformFunc PathTransformFunc
	Transport p2p.Transport
	BootstrapNodes []string
}

type FileServer struct {
	FileServerOpts 

	store *Store
    quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
storeOpts := StoreOpts{
	Root: opts.StorageRoot,
	PathTransformFunc: opts.PathTransformFunc,
}

	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(storeOpts),
		quitch: make(chan struct{}),
	}
} 

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("FileServer stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <- s.Transport.Consume():
			fmt.Printf("FileServer: received message: %v\n", msg)
		case <- s.quitch:
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		go func (addr string) {
		fmt.Println("attemting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Panicln("dial error:", err)
			}
		} (addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	s.loop()

	return nil
}

