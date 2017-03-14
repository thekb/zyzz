package control

import (
	"fmt"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/tcp"
	fb "github.com/google/flatbuffers/go"
	ws "github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/thekb/zyzz/db/models"
	m "github.com/thekb/zyzz/message"
	"time"
)

var (
	PauseHeader   = []byte("p|")
	FrameHeader   = []byte("f|")
	StopHeader    = []byte("s|")
	CommentHeader = []byte("c|")
)

type ControlContext struct {
	WebSocket      *ws.Conn      // pointer to control websocket connection
	currentStream  *Stream       // pointer to current stream
	publisher      bool          // is user tied to the control context publisher
	UserId         int           // user id of the user tied to control context
	loopBack       chan []byte   // for sending messages directly back to client
	publish        bool          // is the user publishing to stream
	pushSocket     mangos.Socket // socket for pushing messages
	subSocket      mangos.Socket // socket for subscribing messages
	closeSubSocket chan bool     // channel for closing sub socket
	streamStarted  bool          // if stream is active on current control context
	builder        *fb.Builder   // flat buffer builder for context
}

// closes control context
func (ctx *ControlContext) Close() {
	ctx.closeSubSocket <- true
	close(ctx.closeSubSocket)
	close(ctx.loopBack)
	ctx.pushSocket.Close()
}

func (ctx *ControlContext) Init(conn *ws.Conn, userId int) {
	ctx.WebSocket = conn
	ctx.UserId = userId
	ctx.builder = fb.NewBuilder(0)
	ctx.closeSubSocket = make(chan bool)
	ctx.loopBack = make(chan []byte, 1)
}

// setup new push socket for current stream
func (ctx *ControlContext) SetupPushSocket() error {
	var err error
	if ctx.pushSocket, err = push.NewSocket(); err != nil {
		fmt.Println("unable to get new push socket:", err)
		return err
	}
	ctx.pushSocket.AddTransport(tcp.NewTransport())
	if err = ctx.pushSocket.Dial(ctx.currentStream.PullSockURL); err != nil {
		fmt.Println("unable to dial to push socket:", err)
		return err
	}
	return nil
}

// setup new sub socket for current stream
func (ctx *ControlContext) SetupSubSocket() error {
	var err error

	if ctx.subSocket, err = sub.NewSocket(); err != nil {
		fmt.Println("unable to get new sub socket:", err)
		return err
	}

	ctx.subSocket.AddTransport(tcp.NewTransport())
	err = ctx.subSocket.Dial(ctx.currentStream.PublishSockURL)
	if err != nil {
		fmt.Println("unable to dial to sub socket:", err)
		return err
	}
	// set receive deadline to 50 ms
	err = ctx.subSocket.SetOption(mangos.OptionRecvDeadline, time.Millisecond * 50)
	if err != nil {
		fmt.Println("unable to set recv deadline:", err)
	}
	if ctx.publish {
		// publisher will subscribe to only stream comments
		err = ctx.subSocket.SetOption(mangos.OptionSubscribe, []byte("c"))
		if err != nil {
			fmt.Println("unable to set subscribe option:", err)
			return err
		}

	} else {
		// subscribe will subscribe to all messages
		err = ctx.subSocket.SetOption(mangos.OptionSubscribe, []byte(""))
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
	defer ctx.subSocket.Close()

	COPY:
	for {
		select {
		// close go routine
		case <- ctx.closeSubSocket:
			fmt.Println("received close to copy")
			break COPY
		case out = <- ctx.loopBack:
			fmt.Println("received message in loopback", out)
			err = ctx.WebSocket.WriteMessage(ws.BinaryMessage, out)
			fmt.Println("sent loopback message on websocket")
			if err != nil {
				if ws.IsCloseError(err) || ws.IsUnexpectedCloseError(err) {
					fmt.Println("websocket connection closed:", err)
					break COPY
				} else {
					fmt.Println("unable to write message to web socket:", err)
				}
			}
		default:
			out, err = ctx.subSocket.Recv()
			if err != nil {
				//fmt.Println("unable to receive from sub socket:", err)
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

func (ctx *ControlContext) UserAllowedToPublish() bool {
	if ctx.currentStream.PublishUser != ctx.UserId {
		fmt.Println("stream publish user, user", ctx.currentStream.PublishUser, ctx.UserId)
		return false
	}
	return true
}

func (ctx *ControlContext) getStreamResponse(streamId, eventId string, err error) []byte {
	ctx.builder.Reset()

	m.ResponseStart(ctx.builder)

	switch err {
	case nil:
		m.ResponseAddStatus(ctx.builder, m.ResponseStatusOK)
	case STREAM_NOT_FOUND:
		m.ResponseAddStatus(ctx.builder, m.ResponseStatusNoStream)
	case STREAM_NOT_ALLOWED:
		m.ResponseAddStatus(ctx.builder, m.ResponseStatusNotAllowed)
	default:
		m.ResponseAddStatus(ctx.builder, m.ResponseStatusOK)
	}

	responseOffset := m.ResponseEnd(ctx.builder)

	streamIdOffset := ctx.builder.CreateString(streamId)
	eventIdOffset := ctx.builder.CreateString(eventId)

	m.StreamMessageStart(ctx.builder)
	m.StreamMessageAddEventId(ctx.builder, eventIdOffset)
	m.StreamMessageAddStreamId(ctx.builder, streamIdOffset)
	m.StreamMessageAddMessageType(ctx.builder, m.MessageResponse)
	m.StreamMessageAddMessage(ctx.builder, responseOffset)
	m.StreamMessageAddTimestamp(ctx.builder, GetCurrentTimeInMilli())
	streamMessageOffset := m.StreamMessageEnd(ctx.builder)

	ctx.builder.Finish(streamMessageOffset)
	return ctx.builder.FinishedBytes()
}

func (ctx *ControlContext) pushMessage(header, in []byte) {
	var msg []byte
	msg = append(msg, header...)
	msg = append(msg, in...)
	err := ctx.pushSocket.Send(msg)
	if err != nil {
		fmt.Println("unable to push message:", err)
	}
}

// send message to client
func (ctx *ControlContext) sendMessageToClient(msg []byte) {
	// if stream has already started send in loopback channel
	if ctx.streamStarted {
		ctx.loopBack <- msg
	} else {
		ctx.WebSocket.WriteMessage(ws.BinaryMessage, msg)
	}
}

// return topic, true if message should be copied to push socket
func (ctx *ControlContext) HandleStreamMessage(db *sqlx.DB, msg []byte) {
	var err error
	var stream *Stream
	streamMessage := m.GetRootAsStreamMessage(msg, 0)
	streamId := string(streamMessage.StreamId())
	eventId := string(streamMessage.EventId())

	stream, err = GetStream(streamId)
	if err != nil {
		fmt.Println("unable to get stream:", err)
		ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))

	}

	//fmt.Println("message lag:", GetCurrentTimeInMilli()-streamMessage.Timestamp())

	switch streamMessage.MessageType() {
	case m.MessageBroadCast:
		fmt.Println("handling stream broadcast")
		// if user is allowed to broadcast on this stream
		if ctx.currentStream == nil || (ctx.currentStream != stream && ctx.UserAllowedToPublish()) {
			ctx.currentStream = stream
			// set publish true
			ctx.publish = true

			// setup push socket
			err = ctx.SetupPushSocket()
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				// send to loopback if stream has started
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))

			}

			err = ctx.SetupSubSocket()
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))
			}
			if ctx.streamStarted {
				// close existing sub socket
				ctx.closeSubSocket <- true
			}
			// start background write subscribe socket to web socket
			ctx.streamStarted = true
			models.SetStreamStatus(db, streamId, models.STATUS_STREAMING)
			go ctx.CopyToWS()
			// send ok back
			ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, nil))
		} else {
			fmt.Println("user not allowed to broadcast")
			ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))

		}
		//TODO what to do with duplicate broadcast messages ?
	case m.MessagePause:
		fmt.Println("handling stream pause")
		if ctx.streamStarted {
			ctx.pushMessage(PauseHeader, msg)
		}
	case m.MessageStop:
		fmt.Println("handling stream stop")
		if ctx.streamStarted {
			models.SetStreamStatus(db, streamId, models.STATUS_STOPPED)
			ctx.pushMessage(StopHeader, msg)
		}
	case m.MessageFrame:
		fmt.Println("handling stream frame")
		if ctx.streamStarted {
			ctx.pushMessage(FrameHeader, msg)
		}
	case m.MessageSubscribe:
		fmt.Println("handling stream subscribe")
		if ctx.currentStream == nil || ctx.currentStream != stream {
			fmt.Println("stream nil or different")
			ctx.currentStream = stream
			ctx.publish = false
			// setup push socket
			err = ctx.SetupPushSocket()
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))
			}
			fmt.Println("push socket setup successfully")
			err = ctx.SetupSubSocket()
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))
			}
			fmt.Println("sub socket setup successfully")
			if ctx.streamStarted {
				// close existing sub socket
				fmt.Println("stream already started, closing exising sub socket")
				ctx.closeSubSocket <- true

			}
			models.IncrementStreamSubscriberCount(db, streamId)
			fmt.Println("incremented subscriber count")
			ctx.streamStarted = true
			go ctx.CopyToWS()
			fmt.Println("started copy to ws goroutine")
			ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, nil))
			ctx.sendMessageToClient(ctx.getStreamStatus(db, eventId, streamId))
			fmt.Println("sent status to client")
		}
	case m.MessageUnSubscribe:
		fmt.Println("handling unsubscribe")
		ctx.currentStream = nil
		ctx.streamStarted = false
		ctx.closeSubSocket <- true
	case m.MessageComment:
		fmt.Println("handling stream comment")
		if ctx.streamStarted {
			ctx.pushMessage(CommentHeader, msg)
		}

	}
}

// sends stream status message to client for current stream
func (ctx *ControlContext) getStreamStatus(db *sqlx.DB, eventId, streamId string) []byte {
	ctx.builder.Reset()
	stream, err := models.GetStreamForShortId(db, streamId)
	if err != nil {
		fmt.Println("unable to get stream:", err)
		return nil
	}

	eventIdOffset := ctx.builder.CreateString(eventId)
	streamIdOffset := ctx.builder.CreateString(streamId)

	m.StatusStart(ctx.builder)
	m.StatusAddStatus(ctx.builder, int8(stream.Status))
	statusOffset := m.StatusEnd(ctx.builder)

	m.StreamMessageStart(ctx.builder)
	m.StreamMessageAddEventId(ctx.builder, eventIdOffset)
	m.StreamMessageAddStreamId(ctx.builder, streamIdOffset)
	m.StreamMessageAddMessageType(ctx.builder, m.MessageStatus)
	m.StreamMessageAddMessage(ctx.builder, statusOffset)
	m.StreamMessageAddTimestamp(ctx.builder, GetCurrentTimeInMilli())
	streamMessageOffset := m.StreamMessageEnd(ctx.builder)
	ctx.builder.Finish(streamMessageOffset)
	return ctx.builder.FinishedBytes()
}
