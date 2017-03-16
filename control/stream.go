package control

import (
	"fmt"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/transport/tcp"
	"sync"
	"github.com/go-mangos/mangos/protocol/pull"
	"errors"
)

const (
	STREAM_TRANSPORT_URL_FORMAT = "tcp://%s:%d"
)

var (
	STREAM_ALREADY_EXISTS = errors.New("Stream Already Exists")
	STREAM_INIT_ERROR = errors.New("Unable to initialize stream")
	STREAM_NOT_FOUND = errors.New("Stream not found")
	STREAM_NOT_ALLOWED = errors.New("Stream not allowed")
)

type Streams struct {
	m       sync.Mutex
	streams map[string]*Stream
}

func (s *Streams) GetStream(streamId string) (*Stream, error) {
	stream, ok := s.streams[streamId]
	if !ok {
		return nil, STREAM_NOT_FOUND
	}
	return stream, nil
}

func (s *Streams) CreateStream(streamId string) error {
	s.m.Lock()
	defer s.m.Unlock()

	stream, ok := s.streams[streamId]
	if !ok {
		fmt.Println("existing stream not found")
		// stream not found
		stream = &Stream{StreamId:streamId}
		err := stream.Init()
		s.streams[streamId] = stream
		// setup required sockets for stream
		if err != nil {
			fmt.Println("Unable to initilaze stream:", err)
			return STREAM_INIT_ERROR
		}
	} else {
		return STREAM_ALREADY_EXISTS
	}
	return nil
}

type Stream struct {
	StreamId       string
	pubSock        mangos.Socket // used to publish data to clients
	pullSock       mangos.Socket // used to pull data from clients
	PublishSockURL string        // url for publish socket
	PullSockURL    string        // url for pull socket
	PublishUser    int           // userid for user who is publishing
}

func (s *Stream)Init() error {
	var err error
	// setup publish socket
	s.pubSock, err = pub.NewSocket()
	if err != nil {
		fmt.Println("unable to create new pub socket:", err)
		return err
	}
	s.pubSock.AddTransport(tcp.NewTransport())
	PubOpenPort, err := GetFreePort()
	if err != nil {
		fmt.Println("unable to get open tcp port:", err)
		return err
	}
	s.PublishSockURL = fmt.Sprintf(STREAM_TRANSPORT_URL_FORMAT, "127.0.0.1", PubOpenPort)
	// timeout sending a message after 10 milliseconds
	err = s.pubSock.Listen(s.PublishSockURL)
	if err != nil {
		fmt.Println("unable to listen at stream transport url:", err)
		return err
	}
	// setup pull socket
	s.pullSock, err = pull.NewSocket()
	if err != nil {
		fmt.Println("unable to create new pull socket")
		return err
	}
	s.pullSock.AddTransport(tcp.NewTransport())
	pullOpenPort, err := GetFreePort()
	s.PullSockURL = fmt.Sprintf(STREAM_TRANSPORT_URL_FORMAT, "127.0.0.1", pullOpenPort)
	err = s.pullSock.Listen(s.PullSockURL)
	if err != nil {
		fmt.Println("unable to listen at stream transport url:", err)
		return err
	}
	// copy in the background
	go s.copy()
	return nil
}

// copies data from pull socket to publish socket
func (s *Stream) copy() {
	var err error
	var msg []byte
	for {
		msg, err = s.pullSock.Recv()
		if err != nil {
			fmt.Println("unable to receive message from pull:", err)
			continue
		}
		err = s.pubSock.Send(msg)
		if err != nil {
			fmt.Println("unable to send message to pub:, err")
		}
	}
}

var StreamMap Streams

func init() {
	StreamMap = Streams{streams:make(map[string]*Stream)}
}