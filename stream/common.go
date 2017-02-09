package stream

import (
	"sync"
	"github.com/ventu-io/go-shortid"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"io"
	"time"
	"errors"
	"net/http"
	"github.com/go-mangos/mangos/transport/ipc"
	"fmt"
)

const (
	DEFAULT_ID_CHARACTERS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ID_SEED = 1729
	IPC_URL_FORMAT = "ipc:///tmp/stream_%s.ipc"
)

var  STREAM_NOT_FOUND = errors.New("No Stream found with id")



var (
	streams Streams
	shortIdGenerator *shortid.Shortid

)

func init() {
	shortIdGenerator, _ = shortid.New(1, shortid.DefaultABC, 2398)
	streams.s = make(map[string]Stream)
}


type Stream struct {
	Name string `json:"name"`
	Id string `json:"id"`
	TransportType string `json:"_"`
	TransportURL string `json:"_"`
}

type Streams struct {
	s map[string]Stream
	m sync.Mutex

}


func GetStream(id string) {

}

// creates stream
func CreateStream(name string) string {
	streams.m.Lock()
	defer streams.m.Unlock()
	id, _ := shortIdGenerator.Generate()
	stream := Stream{}
	stream.Name = name
	stream.TransportURL = fmt.Sprintf(IPC_URL_FORMAT, id)
	streams.s[id] = stream
	return id
}

func PublishStream(id string, reader io.Reader) error {
	if stream, ok := streams.s[id]; ok {
		var sock mangos.Socket
		var err error
		if sock, err = pub.NewSocket(); err != nil {
			return err
		}
		sock.AddTransport(ipc.NewTransport())
		//sock.AddTransport(inproc.NewTransport())
		if err = sock.Listen(stream.TransportURL); err != nil {
			return err
		}

		for {
			fragment := make([]byte, 1024)
			_, err = reader.Read(fragment)
			if err == io.EOF {
				break
			}
			sock.Send(fragment)
			// sleep for 10 milliseconds before reading the next fragment
			time.Sleep(10 *time.Millisecond)
		}
		return nil

	} else {
		return STREAM_NOT_FOUND
	}


}

func SubscribeStream(id string, w http.ResponseWriter) error {
	if stream, ok := streams.s[id]; ok {
		var sock mangos.Socket
		var err error
		var fragment []byte
		flusher, _ := w.(http.Flusher)
		if sock, err = sub.NewSocket(); err != nil {
			return err
		}
		sock.AddTransport(ipc.NewTransport())
		if err = sock.Dial(stream.TransportURL); err != nil {
			return err
		}

		if err = sock.SetOption(mangos.OptionSubscribe, []byte("")); err != nil {
			return err
		}
		for {
			if fragment, err = sock.Recv(); err != nil {
				return err
			}
			w.Write(fragment)
			flusher.Flush()
			// sleep for 10 milliseconds
			time.Sleep(10 * time.Millisecond)
		}
		return nil

	} else {
		return STREAM_NOT_FOUND
	}

}

func DeleteStream(id string) error {
	return nil
}