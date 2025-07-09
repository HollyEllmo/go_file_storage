package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct {}

func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct {}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	buf := make([]byte, 1028)

	// fmt.Println("DefaultDecoder: waiting for data...")
	n, err := r.Read(buf)
	if err != nil {
		fmt.Printf("DefaultDecoder: read error: %v\n", err)
		return  err
	}

	fmt.Printf("DefaultDecoder: received %d bytes: %s\n", n, string(buf[:n]))
	msg.Payload = buf[:n]
   

	return nil
}