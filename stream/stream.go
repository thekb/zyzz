package stream

import (
	"net/http"
	"fmt"
	"github.com/thekb/zyzz/api"
	"github.com/thekb/zyzz/db/models"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"errors"
	"github.com/gorilla/websocket"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/opus"
	"io"
	"github.com/thekb/zyzz/encode"
	"math/rand"
)

const (
	SAMPLE_RATE = 24000
	CHANNELS = 1
	CONTENT_TYPE_AUDIO_WAV = "audio/x-wav;codec=pcm"
	CONTENT_TYPE_AAC = "audio/aac"
	HEADER_CONTENT_TYPE_OPTIONS = "X-Content-Type-Options"
	OPTION_NO_SNIFF = "nosniff"
	CONTENT_TYPE_OPUS = "audio/ogg;codec=opus"
)

var upgrader = websocket.Upgrader{CheckOrigin:CheckOrigin} // use default options

func CheckOrigin(r *http.Request) bool {
	return true
}

type PublishStream struct {
	api.Common
}

func (ps *PublishStream) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(api.SHORT_ID)
	stream, err := models.GetStreamForShortId(ps.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, api.Response{Error:err.Error()})
		return
	}
	if stream.Status != models.STATUS_CREATED {
		ctx.JSON(iris.StatusBadRequest, api.Response{Error:"Publish in progress"})
		return
	}
	var conn *websocket.Conn
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	conn, err = upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to upgrade websocket:", err)
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
	var input []byte
	if sock, err = pub.NewSocket(); err != nil {
		return err
	}
	sock.AddTransport(ipc.NewTransport())

	if err = sock.Listen(stream.TransportUrl); err != nil {
		return err
	}
	//var outputSize int
	var opusEncoder opus.Encoder
	err = opusEncoder.Init(SAMPLE_RATE, CHANNELS, opus.AppAudio)
	if err != nil {
		fmt.Println("unable to init ops encoder:", err)
		return err
	}
	//output := make([]byte, 1024)
	//encoderInput := make([]int16, 480)
	//TODO optimize with io reader
	for {
		_, input, err = conn.ReadMessage()
		if err != nil {
			fmt.Println("error reading message:", err)
			break
		}
		/*
		err = binary.Read(bytes.NewReader(input), binary.LittleEndian, &encoderInput)
		if err != nil {
			fmt.Println("error reading message:", err)
			continue
		}
		//outputSize, err = opusEncoder.Encode(encoderInput, output)
		if err != nil {
			fmt.Println("unable to encode pcm to opus:", err)
			continue
		}

		if outputSize > 2 {
			sock.Send(output[:outputSize])
		}
		*/
		sock.Send(input)


	}
	conn.Close()
	return errors.New("Stream Closed")
}



type SubscribeStream struct {
	api.Common
}

func (ss *SubscribeStream) Serve(ctx *iris.Context){
	shortId := ctx.GetString(api.SHORT_ID)
	var err error
	stream, err := models.GetStreamForShortId(ss.DB, shortId)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &api.Response{Error:err.Error()})
		return
	}
	if stream.Status != models.STATUS_STREAMING {
		ctx.JSON(iris.StatusBadRequest, &api.Response{Error:"No active stream"})
		return
	}
	models.IncrementStreamSubscriberCount(ss.DB, shortId)
	ctx.SetHeader(api.HEADER_CONTENT_TYPE, CONTENT_TYPE_OPUS)
	ctx.SetHeader(HEADER_CONTENT_TYPE_OPTIONS, OPTION_NO_SNIFF)
	comments := make(map[string]string)
	comments["NAME"] = "TEST STREAM"
	comments["ALBUM"] = "TEST ALBUM"
	opusOggStream := encode.OpusOggStream{
		StreamId: rand.Int31(),
		Channels: 1,
		PreSkip: 0,
		InputSampleRate: 24000,
		OutPutGain: 0,
		ChannelMap: 0,
		VendorString: "thekb zyzz encoder",
		Comments: comments,
		FrameSize: 10.0,
	}
	// write ogg/opus headers to stream
	ctx.StreamWriter(func (w io.Writer) bool{
		w.Write(opusOggStream.Start())
		return false
	})
	var sock mangos.Socket
	var fragment []byte
	if sock, err = sub.NewSocket(); err != nil {
		ctx.JSON(iris.StatusInternalServerError, &api.Response{Error:err.Error()})
		return
	}
	sock.AddTransport(ipc.NewTransport())
	if err = sock.Dial(stream.TransportUrl); err != nil {
		ctx.JSON(iris.StatusInternalServerError, &api.Response{Error:err.Error()})
		return
	}

	if err = sock.SetOption(mangos.OptionSubscribe, []byte("")); err != nil {
		ctx.JSON(iris.StatusInternalServerError, &api.Response{Error:err.Error()})
		return
	}
	for {
		if fragment, err = sock.Recv(); err != nil {
			ctx.JSON(iris.StatusInternalServerError, &api.Response{Error:err.Error()})
			return
		}
		ctx.StreamWriter(func(w io.Writer) bool {
			w.Write(opusOggStream.FlushPacket(fragment))
			return false
		})

	}
}

type WebSocketSubscriber struct {
	api.Common
}

func (wss *WebSocketSubscriber) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(api.SHORT_ID)
	var err error
	stream, err := models.GetStreamForShortId(wss.DB, shortId)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &api.Response{Error:err.Error()})
		return
	}
	if stream.Status != models.STATUS_STREAMING {
		ctx.JSON(iris.StatusBadRequest, &api.Response{Error:"No active stream"})
		return
	}
	var conn *websocket.Conn
	conn, err = upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to upgrade websocket:", err)
		return
	}

	models.IncrementStreamSubscriberCount(wss.DB, shortId)

	/*
	comments := make(map[string]string)
	comments["NAME"] = "TEST STREAM"
	comments["ALBUM"] = "TEST ALBUM"
	opusOggStream := encode.OpusOggStream{
		StreamId: rand.Int31(),
		Channels: 1,
		PreSkip: 0,
		InputSampleRate: 24000,
		OutPutGain: 0,
		ChannelMap: 0,
		VendorString: "thekb zyzz encoder",
		Comments: comments,
		FrameSize: 10.0,
	}
	err = conn.WriteMessage(websocket.BinaryMessage, opusOggStream.Start())
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to write stream headers, closing:", err)
		conn.Close()
		return
	}
	*/
	var sock mangos.Socket
	var fragment []byte
	if sock, err = sub.NewSocket(); err != nil {
		conn.WriteJSON(&api.Response{Error:err.Error()})
		return
	}
	sock.AddTransport(ipc.NewTransport())
	if err = sock.Dial(stream.TransportUrl); err != nil {
		conn.WriteJSON(&api.Response{Error:err.Error()})
		return
	}

	if err = sock.SetOption(mangos.OptionSubscribe, []byte("")); err != nil {
		conn.WriteJSON(&api.Response{Error:err.Error()})
		return
	}
	for {
		if fragment, err = sock.Recv(); err != nil {
			ctx.Log(iris.ProdMode, "unabel to receive from socket:", err)
			continue
		}
		conn.WriteMessage(websocket.BinaryMessage, fragment)
	}

}