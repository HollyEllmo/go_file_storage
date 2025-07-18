package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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

	peerLock sync.Mutex
	peers map[string]p2p.Peer
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
		peers: make(map[string]p2p.Peer),
	}
} 


func (s *FileServer) broadcast(msg *Message) error {
peers := []io.Writer{}

for _, peer := range s.peers {
	peers = append(peers, peer)
}

mw := io.MultiWriter(peers...)
return gob.NewEncoder(mw).Encode(msg)
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key string
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to disk
	// 2. broadcast this file to all known peers in the network

	buf := new(bytes.Buffer)
	msg := Message{
		Payload: MessageStoreFile{
			Key: key,
		},
	}
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	fmt.Println("StoreData: Sent storagekey, now sleeping 3 seconds...")
	time.Sleep(time.Second * 3)
	fmt.Println("StoreData: Sleep done, sending file...")

	payload := []byte("THIS LARGE FILE")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}
	fmt.Println("StoreData: File sent!")

	return nil

//    buf := new(bytes.Buffer)
//    tee := io.TeeReader(r, buf)

//    if err := s.store.Write(key, tee); err != nil {
// 	return err
//   }

//    p := &DataMessage{
// 	Key:  key,
// 	Data: buf.Bytes(),
//    }

// 	return s.broadcast(&Message{
// 		From: "todo",
// 		Payload: p,
// 	})
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s", p.RemoteAddr().String())

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("FileServer stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <- s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Panicln(err)
				return
			}
		if err := s.handleMessage(rpc.From, &msg); err != nil {
			log.Panicln(err)
			return
		}

       
		case <- s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)
	}

	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer %s not found in peers map", from)
	}

	return s.store.Write(msg.Key, peer)
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
	}
	for _, addr := range s.BootstrapNodes {
		go func (addr string) {
			fmt.Println("attempting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Printf("failed to dial %s: %v\n", addr, err)
				return
			}
			fmt.Printf("successfully connected to %s\n", addr)
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

func init() {
	gob.Register(MessageStoreFile{})
}

