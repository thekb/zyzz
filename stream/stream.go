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
	"errors"
	"github.com/gorilla/websocket"
)

const (
	CONTENT_TYPE_AUDIO_WAV = "audio/x-wav;codec=pcm"
	CONTENT_TYPE_AAC = "audio/aac"
	HEADER_CONTENT_TYPE_OPTIONS = "X-Content-Type-Options"
	OPTION_NO_SNIFF = "nosniff"
)

var upgrader = websocket.Upgrader{CheckOrigin:CheckOrigin} // use default options

func CheckOrigin(r *http.Request) bool {
	return true
}

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
	if stream.Status != models.STATUS_CREATED {
		ps.SendErrorJSON(rw, "Publish in progress", http.StatusBadRequest)
		return
	}
	var conn *websocket.Conn
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	conn, err = upgrader.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Println("unable to upgrade to websocket:", err)
		return
	}


	models.SetStreamStatus(ps.DB, shortId, models.STATUS_STREAMING)
	ps.publish(stream, conn)
	// when the publish loop exits, set set stream status
	models.SetStreamStatus(ps.DB, shortId, models.STATUS_STOPPED)
}

func (ps *PublishStream) publish(stream models.Stream, conn *websocket.Conn) error {
	var sock mangos.Socket
	var err error
	var fragment []byte
	if sock, err = pub.NewSocket(); err != nil {
		return err
	}
	sock.AddTransport(ipc.NewTransport())

	if err = sock.Listen(stream.TransportUrl); err != nil {
		return err
	}
	var encoder *fdkaac.AacEncoder
	encoder = GetNewEncoder()
	for {
		_, fragment, err = conn.ReadMessage()
		fmt.Println(fragment)
		if err != nil {
			break
		}

		// if encoder is present encode and publish
		var aacFragment []byte
		aacFragment, err = encoder.Encode(fragment)
		if err != nil {
			fmt.Println("unable to encode pcm to aac:", err)
		}
		sock.Send(aacFragment)

	}
	conn.Close()
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
	if stream.Status != models.STATUS_STREAMING {
		ss.SendErrorJSON(rw, "No Active Stream", http.StatusBadRequest)
		return
	}
	models.IncrementStreamSubscriberCount(ss.DB, shortId)
	rw.Header().Set(api.HEADER_CONTENT_TYPE, CONTENT_TYPE_AAC)
	rw.Header().Set(HEADER_CONTENT_TYPE_OPTIONS, OPTION_NO_SNIFF)
	ss.subscribe(stream, rw)
	return
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