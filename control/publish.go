package control

import (
	"fmt"
	"sync"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/transport/ipc"
)

type StreamControl struct {
	m       sync.Mutex
	socket  map[string]mangos.Socket
	control map[string]chan []byte // channel for control messages
	comment map[string]chan []byte // channel for comments
}


// creates publish socket for stream
func (sc *StreamControl) InitStream(streamId, transportURL string) error {
	sc.m.Lock()
	fmt.Printf("initializing stream control for stream %s, transport %s", streamId, transportURL)
	defer sc.m.Unlock()
	var socket mangos.Socket
	var err error
	if socket, err = pub.NewSocket(); err != nil {
		return err
	}
	socket.AddTransport(ipc.NewTransport())
	if err = socket.Listen(transportURL); err != nil {
		return err
	}
	sc.socket[streamId] = socket
	sc.control[streamId] = make(chan []byte, 10)
	sc.comment[streamId] = make(chan []byte, 10)
	// separate comments from controls/frames
	go sendToSocket(sc.control[streamId], sc.socket[streamId])
	go sendToSocket(sc.comment[streamId], sc.socket[streamId])
	return nil
}

// sends messages received from channel to socket
// for serializing publish on a stream
func sendToSocket(receive chan []byte, socket mangos.Socket) {
	var msg []byte
	var err error
	for {
		msg = <- receive
		err = socket.Send(msg)
		if err != nil {
			fmt.Println("unable to send message to socket:", err)
		}
	}
}



