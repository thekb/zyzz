package control

import (
	"gopkg.in/kataras/iris.v6"
	ws "github.com/gorilla/websocket"
	"net/http"
	"fmt"
	m "github.com/thekb/zyzz/message"
	fb "github.com/google/flatbuffers/go"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/jmoiron/sqlx"
	"time"
	"github.com/thekb/zyzz/db/models"
)

type Control struct {
	DB *sqlx.DB
}

type ControlContext struct {
	WebSocket *ws.Conn // pointer to control websocket connection
	CurrentStream *Stream // pointer to current stream
	Publisher bool // is user tied to the control context publisher
	UserId int // user id of the user tied to control context
	Loopback chan []byte // for sending messages directly back to client
	Publish bool // is the user publishing to stream
	PushSocket mangos.Socket // socket for pushing messages
	SubSocket mangos.Socket // socket for subscribing messages
	CloseSubSocket chan bool // channel for closing sub socket
	StreamStarted bool // if stream is active on current control context
	Builder *fb.Builder
}

// setup new push socket for current stream
func (ctx *ControlContext) SetupPushSocket() error {
	var err error
	if ctx.PushSocket, err = push.NewSocket(); err != nil {
		fmt.Println("unable to get new push socket:", err)
		return err
	}
	ctx.PushSocket.AddTransport(tcp.NewTransport())
	if err = ctx.PushSocket.Dial(ctx.CurrentStream.PullSockURL); err != nil {
		fmt.Println("unable to dial to push socket:", err)
		return err
	}
	return nil
}

// setup new sub socket for current stream
func (ctx *ControlContext) SetupSubSocket() error {
	var err error

	if ctx.SubSocket, err = sub.NewSocket(); err != nil {
		fmt.Println("unable to get new sub socket:", err)
		return err
	}

	ctx.SubSocket.AddTransport(tcp.NewTransport())
	err = ctx.SubSocket.Dial(ctx.CurrentStream.PublishSockURL)
	if err != nil {
		fmt.Println("unable to dial to sub socket:", err)
		return err
	}
	if ctx.Publish {
		// publisher will subscribe to only stream comments
		err = ctx.SubSocket.SetOption(mangos.OptionSubscribe, []byte("c"))
		if err != nil {
			fmt.Println("unable to set subscribe option:", err)
			return err
		}

	} else {
		// subscribe will subscribe to all messages
		err = ctx.SubSocket.SetOption(mangos.OptionSubscribe, []byte(""))
		if err != nil {
			fmt.Println("unable to set subscribe option:", err)
			return err
		}
	}
	return nil
}

// read from sub socket and writes to websocket
func (ctx *ControlContext) CopyToWS() {
	var out []byte
	var err error
	defer ctx.SubSocket.Close()
	COPY:
	for {
		select {
		// close go routine
		case <- ctx.CloseSubSocket:
			break
		case out = <- ctx.Loopback:
			err = ctx.WebSocket.WriteMessage(ws.BinaryMessage, out)
			if err != nil {
				if ws.IsCloseError(err) || ws.IsUnexpectedCloseError(err) {
					fmt.Println("websocket connection closed:", err)
					break COPY
				} else {
					fmt.Println("unable to write message to web socket:", err)
				}
			}
		default:
			out, err = ctx.SubSocket.Recv()
			if err != nil {
				fmt.Println("unable to receive from sub socket:", err)
				continue
			}
			// first 2 bytes contain topic and delimiter
			err = ctx.WebSocket.WriteMessage(ws.BinaryMessage, out[2:])
			if err != nil {
				if ws.IsCloseError(err) || ws.IsUnexpectedCloseError(err) {
					fmt.Println("websocket connection closed:", err)
					break COPY
				} else {
					fmt.Println("unable to write message to web socket:", err)
				}
			}

		}
	}

}


var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func (c *Control) Serve(ctx *iris.Context) {
	// verify if user is authenticated etc etc
	var err error
	var wsc *ws.Conn

	var userId int

	// validate user session
	/*
	userId, err = strconv.ParseInt(ctx.RequestHeader("X-User-Id"), 10, 0)
	if err != nil {
		fmt.Println("unable to get user id from header:", err)
		return
	}
	*/



	userId, err = ctx.Session().GetInt("id")
	if err != nil {
		fmt.Println("unable to get user id:", err)
		//ctx.JSON(iris.StatusBadRequest, &api.Response{Error: "user not authenticated"})
		return
	}

	

	// upgrade current control socket get request to websocket
	wsc, err = upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, ctx.ResponseWriter.Header())
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, map[string]string{"Error":err.Error()})
		return
	}
	// upgrade to websocket end

	// once the upgrade is successful create control context
	controlContext := &ControlContext{
		WebSocket:wsc,
		UserId:userId,
		Builder: fb.NewBuilder(0),
		CloseSubSocket: make(chan bool),
		Loopback: make(chan []byte),
	}

	var in []byte
	var copy bool
	var header []byte
	// read from websocket and push to stream
	for {
		fmt.Println("reading message")
		_, in, err = wsc.ReadMessage()
		if err != nil  {
			// if websocket connection is closed break out of the read loop
			if ws.IsUnexpectedCloseError(err) || ws.IsCloseError(err) {
				fmt.Println("websocket closed:", err)
				break
			}
			// else continue
			fmt.Println("unable to read from websocket:", err)
			continue
		}

		header, copy = controlContext.HandleStreamMessage(c.DB, in)

		if copy {
			var msg []byte
			msg = append(msg, header...)
			msg = append(msg, in...)
			err = controlContext.PushSocket.Send(msg)
			if err != nil {
				fmt.Println("unable to send to push socket:", err)
			}
		}

	}


}

func (ctx *ControlContext) UserAllowedToPublish() bool {
	if ctx.CurrentStream.PublishUser != ctx.UserId {
		fmt.Println("stream publish user, user", ctx.CurrentStream.PublishUser, ctx.UserId)
		return false
	}
	return true
}

func (ctx *ControlContext) GetStreamResponse(streamId, eventId string, err error) []byte {
	ctx.Builder.Reset()

	m.StreamResponseStart(ctx.Builder)

	switch err {
	case nil:
		m.StreamResponseAddStatus(ctx.Builder, m.StatusOK)
	case STREAM_NOT_FOUND:
		m.StreamResponseAddStatus(ctx.Builder, m.StatusNoStream)
	case STREAM_NOT_ALLOWED:
		m.StreamResponseAddStatus(ctx.Builder, m.StatusNotAllowed)
	default:
		m.StreamResponseAddStatus(ctx.Builder, m.StatusError)
	}

	responseOffset := m.StreamResponseEnd(ctx.Builder)

	streamIdOffset := ctx.Builder.CreateString(streamId)
	eventIdOffset := ctx.Builder.CreateString(eventId)

	m.StreamMessageStart(ctx.Builder)
	m.StreamMessageAddEventId(ctx.Builder, eventIdOffset)
	m.StreamMessageAddStreamId(ctx.Builder, streamIdOffset)
	m.StreamMessageAddMessageType(ctx.Builder, m.MessageStreamResponse)
	m.StreamMessageAddMessage(ctx.Builder, responseOffset)
	m.StreamMessageAddTimestamp(ctx.Builder, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(ctx.Builder)

	ctx.Builder.Finish(streamMessageOffset)
	return ctx.Builder.FinishedBytes()
}


// send error message to client
func (ctx *ControlContext) SendErrorMessage(streamId, eventId string, err error) {
	// if stream has already started send in loopback channel
	if ctx.StreamStarted {
		ctx.Loopback <- ctx.GetStreamResponse(streamId, eventId, err)
	} else {
		ctx.WebSocket.WriteMessage(ws.BinaryMessage, ctx.GetStreamResponse(streamId, eventId, err))
	}
}

// return topic, true if message should be copied to push socket
func (ctx *ControlContext) HandleStreamMessage(db *sqlx.DB ,msg []byte) ([]byte, bool) {
	var err error
	var stream *Stream
	streamMessage := m.GetRootAsStreamMessage(msg, 0)
	streamId := string(streamMessage.StreamId())
	eventId := string(streamMessage.EventId())

	stream, err = GetStream(streamId)
	if err != nil {
		fmt.Println("unable to get stream:", err)
		ctx.SendErrorMessage(streamId, eventId, err)
		return nil, false

	}

	fmt.Println("message lag:", GetCurrentTimeInMilli() - streamMessage.Timestamp())

	switch streamMessage.MessageType() {
	case m.MessageStreamBroadCast:
		fmt.Println("handling stream broadcast")
		// if user is allowed to broadcast on this stream
		if ctx.CurrentStream == nil || (ctx.CurrentStream != stream && ctx.UserAllowedToPublish()) {
			ctx.CurrentStream = stream
			// set publish true
			ctx.Publish = true

			// setup push socket
			err = ctx.SetupPushSocket()
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				// send to loopback if stream has started
				ctx.SendErrorMessage(streamId, eventId, err)

			}

			err = ctx.SetupSubSocket()
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				ctx.SendErrorMessage(streamId, eventId, err)
			}
			if ctx.StreamStarted {
				// close existing sub socket
				ctx.CloseSubSocket <- true
			}
			// start background write subscribe socket to web socket
			ctx.StreamStarted = true
			models.SetStreamStatus(db, streamId, models.STATUS_STREAMING)
			go ctx.CopyToWS()
			// send ok back
			ctx.SendErrorMessage(streamId, eventId, nil)
		} else {
			fmt.Println("user not allowed to broadcast")
			ctx.SendErrorMessage(streamId, eventId, err)

		}
		//TODO what to do with duplicate broadcast messages ?
		return nil, false
	case m.MessageStreamPause:
		fmt.Println("handling stream pause")
		return []byte("p|"), true
	case m.MessageStreamStop:
		fmt.Println("handling stream stop")
		models.SetStreamStatus(db, streamId, models.STATUS_STOPPED)
		return []byte("s|"), true
	case m.MessageStreamFrame:
		fmt.Println("handling stream frame")
		return []byte("f|"), true
	case m.MessageStreamSubscribe:
		if ctx.CurrentStream != stream {
			ctx.CurrentStream = stream
			ctx.Publish = false
			// setup push socket
			err = ctx.SetupPushSocket()
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				ctx.SendErrorMessage(streamId, eventId, err)
			}

			err = ctx.SetupSubSocket()
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				ctx.SendErrorMessage(streamId, eventId, err)
			}

			if ctx.StreamStarted {
				// close existing sub socket
				ctx.CloseSubSocket <- true

			}
			models.IncrementStreamSubscriberCount(db, streamId)
			go ctx.CopyToWS()
		}
		return nil, false
	case m.MessageStreamComment:
		fmt.Println("handling stream comment")
		return []byte("c|"), true
	}
	return nil, false
}

func GetTimeInMillis() int64 {
	 return time.Now().UnixNano() / int64(time.Millisecond)
}
