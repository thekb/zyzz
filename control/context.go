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
	PauseHeader = []byte("p|")
	FrameHeader = []byte("f|")
	StopHeader = []byte("s|")
	CommentHeader = []byte("c|")
	ActiveListenerHeader = []byte("a|")
)

type ControlContext struct {
	WebSocket *ws.Conn      // pointer to control websocket connection
	stream    *Stream       // pointer to current stream
	publisher bool          // is user tied to the control context publisher
	UserId    int           // user id of the user tied to control context
	loopBack  chan []byte   // for sending messages directly back to client
	publish   bool          // is the user publishing to stream
	push      mangos.Socket // socket for pushing messages
	sub       mangos.Socket // socket for subscribing messages
	closeCopy chan bool     // channel for closing sub socket
	active    bool          // if stream is active on current control context
	builder   *fb.Builder   // flat buffer builder for context
}

// closes control context
func (ctx *ControlContext) Close() {
	ctx.closeCopy <- true
	close(ctx.closeCopy)
	close(ctx.loopBack)
	ctx.push.Close()
}

func (ctx *ControlContext) Init(conn *ws.Conn, userId int) {
	ctx.WebSocket = conn
	ctx.UserId = userId
	ctx.builder = fb.NewBuilder(0)
	ctx.closeCopy = make(chan bool)
	ctx.loopBack = make(chan []byte, 1)
}

// setup new push socket for current stream
func (ctx *ControlContext) SetupPushSocket() error {
	var err error
	if ctx.push, err = push.NewSocket(); err != nil {
		fmt.Println("unable to get new push socket:", err)
		return err
	}
	ctx.push.AddTransport(tcp.NewTransport())
	if err = ctx.push.Dial(ctx.stream.PullSockURL); err != nil {
		fmt.Println("unable to dial to push socket:", err)
		return err
	}
	return nil
}

// setup new sub socket for current stream
func (ctx *ControlContext) SetupSubSocket() error {
	var err error

	if ctx.sub, err = sub.NewSocket(); err != nil {
		fmt.Println("unable to get new sub socket:", err)
		return err
	}

	ctx.sub.AddTransport(tcp.NewTransport())
	err = ctx.sub.Dial(ctx.stream.PublishSockURL)
	if err != nil {
		fmt.Println("unable to dial to sub socket:", err)
		return err
	}
	// set receive deadline to 10 ms
	err = ctx.sub.SetOption(mangos.OptionRecvDeadline, time.Millisecond * 10)
	if err != nil {
		fmt.Println("unable to set recv deadline:", err)
	}
	if ctx.publish {
		// publisher will subscribe to only stream comments
		err = ctx.sub.SetOption(mangos.OptionSubscribe, []byte("c"))
		err = ctx.sub.SetOption(mangos.OptionSubscribe, []byte("a"))
		if err != nil {
			fmt.Println("unable to set subscribe option:", err)
			return err
		}

	} else {
		// subscribe will subscribe to all messages
		err = ctx.sub.SetOption(mangos.OptionSubscribe, []byte(""))
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
	defer ctx.sub.Close()

	COPY:
	for {
		select {
		// close go routine
		case <-ctx.closeCopy:
			fmt.Println("received close to copy")
			break COPY
		case out = <-ctx.loopBack:
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
			out, err = ctx.sub.Recv()
			if err != nil {
				if err != mangos.ErrRecvTimeout {
					fmt.Println("unable to receive from sub socket:", err)
				}
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
	if ctx.stream.PublishUser != ctx.UserId {
		fmt.Println("stream publish user, user", ctx.stream.PublishUser, ctx.UserId)
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
		m.ResponseAddStatus(ctx.builder, m.ResponseStatusError)
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
	err := ctx.push.Send(msg)
	if err != nil {
		fmt.Println("unable to push message:", err)
	}
}

// send message to client
func (ctx *ControlContext) sendMessageToClient(msg []byte) {
	// if stream has already started send in loopback channel
	if ctx.active {
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

	stream, err = StreamMap.GetStream(streamId)
	if err != nil {
		fmt.Println("unable to get stream:", err)
		ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))

	}

	//fmt.Println("message lag:", GetCurrentTimeInMilli()-streamMessage.Timestamp())

	switch streamMessage.MessageType() {
	case m.MessageBroadCast:
		fmt.Println("handling stream broadcast")
		// if user is allowed to broadcast on this stream
		if ctx.stream == nil {
			fmt.Println("no active stream")
			ctx.active = false
			// set stream
			ctx.stream = stream
			// set publish true
			if !ctx.UserAllowedToPublish() {
				fmt.Println("user not allowed to broadcast")
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, STREAM_NOT_ALLOWED))
				return
			}
			ctx.publish = true
			// setup push socket
			err = ctx.SetupPushSocket()
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				// send to loopback if stream has started
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))
				return
			}
			// setup sub socket
			err = ctx.SetupSubSocket()
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, err))
				return
			}
			// update stream status
			models.SetStreamStatus(db, streamId, models.STATUS_STREAMING)
			// send response back to client
			ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, nil))
			// start background copy subscribe socket to web socket
			ctx.active = true
			go ctx.CopyToWS()
		}
	//TODO what to do with duplicate broadcast messages ?
	case m.MessagePause:
		fmt.Println("handling stream pause")
		if ctx.active {
			ctx.pushMessage(PauseHeader, msg)
		}
	case m.MessageStop:
		fmt.Println("handling stream stop")
		if ctx.active {
			fmt.Println("stream already started cleaning up")
			models.SetStreamStatus(db, streamId, models.STATUS_STOPPED)
			// send stop message to subscribers
			ctx.pushMessage(StopHeader, msg)
			ctx.closeCopy <- true
			ctx.push.Close()
			ctx.stream = nil
		}
	// TODO should we cleanup after stop ?
	case m.MessageFrame:
		//fmt.Println("handling stream frame")
		if ctx.active {
			ctx.pushMessage(FrameHeader, msg)
		}
	case m.MessageSubscribe:
		fmt.Println("handling stream subscribe")
		if ctx.stream == nil {
			fmt.Println("stream not active")
			ctx.stream = stream
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
			models.IncrementStreamSubscriberCount(db, streamId)
			models.IncrementActiveListenersCount(db, streamId)
			fmt.Println("incremented subscriber count")
			ctx.sendMessageToClient(ctx.getStreamResponse(streamId, eventId, nil))
			ctx.sendMessageToClient(ctx.getStreamStatus(db, eventId, streamId))
			fmt.Println("sent status to client")
			ctx.active = true
			go ctx.CopyToWS()
			fmt.Println("started copy to ws goroutine")
		}
	case m.MessageUnSubscribe:
		fmt.Println("handling unsubscribe")
		if ctx.active {
			ctx.stream = nil
			ctx.active = false
			ctx.closeCopy <- true
			ctx.push.Close()
		}
		models.DecrementActiveListenersCount(db, streamId)
		actMsg := ctx.GetStreamActilveListenersMessage(db, streamId, eventId)
		ctx.pushMessage(ActiveListenerHeader, actMsg)
	case m.MessageComment:
		//fmt.Println("handling stream comment")
		if ctx.active {
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
	m.StatusAddSubscribeCount(ctx.builder, int32(stream.SubscriberCount))
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

func (ctx *ControlContext) GetStreamActilveListenersMessage(db *sqlx.DB, streamId, eventId string) []byte {
	ctx.builder.Reset()

	stream, err := models.GetStreamForShortId(db, streamId)
	if err != nil {
		fmt.Println("unable to get stream:", err)
		return nil
	}
	eventIdOffset := ctx.builder.CreateString(eventId)
	streamIdOffset := ctx.builder.CreateString(streamId)

	m.ActiveListenersStart(ctx.builder)
	m.ActiveListenersAddActiveListeners(ctx.builder, int32(stream.ActiveListeners))
	streamALOffset := m.ActiveListenersEnd(ctx.builder)

	m.StreamMessageStart(ctx.builder)
	m.StreamMessageAddEventId(ctx.builder, eventIdOffset)
	m.StreamMessageAddStreamId(ctx.builder, streamIdOffset)
	m.StreamMessageAddMessageType(ctx.builder, m.MessageActiveListeners)
	m.StreamMessageAddMessage(ctx.builder, streamALOffset)
	m.StreamMessageAddTimestamp(ctx.builder, GetCurrentTimeInMilli())
	streamMessageOffset := m.StreamMessageEnd(ctx.builder)
	ctx.builder.Finish(streamMessageOffset)
	return ctx.builder.FinishedBytes()
}
