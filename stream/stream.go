package stream

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thekb/zyzz/api"
	"github.com/thekb/zyzz/db/models"
	"github.com/winlinvip/go-fdkaac/fdkaac"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"io"
	"errors"
)

const (
	CONTENT_TYPE_AUDIO_WAV = "audio/x-wav;codec=pcm;rate=44100"
	CONTENT_TYPE_AAC = "audio/aac"
)


type PublishStream struct {
	api.Common
}

func (ps PublishStream) ServeHTTP(rw http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	shortId := vars[api.SHORT_ID]
	stream, err := models.GetStreamForShortId(ps.DB, shortId)
	if err != nil {
		ps.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(api.HEADER_CONTENT_TYPE)
	// if input stream is wav setup encoder
	err = ps.publish(stream, r.Body, contentType)

	ps.SendErrorJSON(rw, err.Error(), http.StatusGone)
}

func (ps *PublishStream) publish(stream models.Stream, r io.Reader, contentType string) error {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		return err
	}
	sock.AddTransport(ipc.NewTransport())

	if err = sock.Listen(stream.TransportUrl); err != nil {
		return err
	}
	var encoder *fdkaac.AacEncoder
	if contentType == CONTENT_TYPE_AUDIO_WAV {
		encoder = GetNewEncoder()
	}
	for {
		fragment := make([]byte, 4096)

		_, err = r.Read(fragment)
		if err == io.EOF {
			break
		}

		// if encoder is present encode and publish
		if encoder != nil {
			var aacFragment []byte
			aacFragment, err = encoder.Encode(fragment)
			if err != nil {
				fmt.Println("unable to encode pcm to aac:", err)
			}
			sock.Send(aacFragment)
		} else {
			sock.Send(fragment)
		}
	}
	return errors.New("Stream Closed")
}



type SubscribeStream struct {
	api.Common
}

func (ss SubscribeStream) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortId := vars[api.SHORT_ID]
	stream, err := models.GetStreamForShortId(ss.DB, shortId)

	if err != nil {
		ss.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "audio/aac")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	err = ss.subscribe(stream, rw)

	ss.SendErrorJSON(rw, err.Error(), http.StatusGone)

}

func (ss *SubscribeStream) subscribe(stream models.Stream, w http.ResponseWriter) error {
	var sock mangos.Socket
	var err error
	var fragment []byte
	flusher, _ := w.(http.Flusher)
	if sock, err = sub.NewSocket(); err != nil {
		return err
	}
	sock.AddTransport(ipc.NewTransport())
	if err = sock.Dial(stream.TransportUrl); err != nil {
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
	}
	return errors.New("Stream Closed")


}