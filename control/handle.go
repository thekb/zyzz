package control

import (
	"gopkg.in/kataras/iris.v6"
	ws "github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"github.com/thekb/zyzz/message"
	fb "github.com/google/flatbuffers/go"
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/go-mangos/mangos/protocol/sub"
	"time"
	"strconv"
)

type Control struct {
	currentStream *Stream
	publisher bool
	userId int
	loopback chan []byte // for sending messages directly back to client
	publish bool
	pushSocket mangos.Socket
	closeSubSocket chan bool // for closing sub socket

}


func (c *Control) SetupPushSocket(stream *Stream) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error
	if sock, err = push.NewSocket(); err != nil {
		fmt.Println("unable to get new push socket:", err)
		return sock, err
	}
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(stream.PullSockURL); err != nil {
		fmt.Println("unable to dial to push socket:", err)
		return sock, err
	}
	return sock, nil
}

// read from sub socket and writes to websocket
func (c *Control) CopyToWS(conn *ws.Conn, subSocket mangos.Socket) {
	var out []byte
	var err error
	defer subSocket.Close()

	for {
		select {
		// close go routine
		case <- c.closeSubSocket:
			break
		case out = <- c.loopback:
			err = conn.WriteMessage(ws.BinaryMessage, out)
			if err != nil {
				if ws.IsCloseError(err) || ws.IsUnexpectedCloseError(err) {
					fmt.Println("websocket connection closed:", err)
					break
				} else {
					fmt.Println("unable to write message to web socket:", err)
				}
			}
		default:
			out, err = subSocket.Recv()
			if err != nil {
				fmt.Println("unable to receive from sub socket:", err)
				continue
			}
			err = conn.WriteMessage(ws.BinaryMessage, out)
			if err != nil {
				if ws.IsCloseError(err) || ws.IsUnexpectedCloseError(err) {
					fmt.Println("websocket connection closed:", err)
					break
				} else {
					fmt.Println("unable to write message to web socket:", err)
				}
			}

		}
	}

}

func (c *Control) SetupSubSocket(stream *Stream) (mangos.Socket, error) {
	var err error
	var sock mangos.Socket

	if sock, err = sub.NewSocket(); err != nil {
		fmt.Println("unable to get new sub socket:", err)
		return sock, err
	}

	sock.AddTransport(tcp.NewTransport())

	err = sock.DialOptions(stream.PublishSockURL, map[string]interface{}{
		mangos.OptionSubscribe: []byte(""),
		// set receive time out to 10 milliseconds
		mangos.OptionRecvDeadline: time.Millisecond * 10,
	})
	if err != nil {
		fmt.Println("unable to dial to sub socket:", err)
		return sock, err
	}
	return sock, nil
}

func (c *Control) Serve(ctx *iris.Context) {
	// verify if user is authenticated etc etc
	var err error
	var wsc *ws.Conn

	var upgrader = ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var userId int64

	userId, err = strconv.ParseInt(ctx.RequestHeader("X-User-Id"), 10, 0)
	if err != nil {
		fmt.Println("unable to get user id from header:", err)
		return
	}
	c.userId = int(userId)

	/*
	c.userId, err = ctx.Session().GetInt("id")
	if err != nil {
		fmt.Println("unable to get user id:", err)
		ctx.JSON(iris.StatusBadRequest, &api.Response{Error: "user not authenticated"})
		return
	}
	*/


	wsc, err = upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, ctx.ResponseWriter.Header())
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, map[string]string{"Error":err.Error()})
		return
	}


	// initialize close sub channel
	c.closeSubSocket = make(chan bool, 1)
	// initialize loop back channel
	c.loopback = make(chan []byte)

	var in []byte
	var copy bool
	// read from websocket and push to stream
	for {
		_, in, err = wsc.ReadMessage()
		if err != nil  {
			// if websocket connection is closed
			if ws.IsUnexpectedCloseError(err) || ws.IsCloseError(err) {
				fmt.Println("websocket closed:", err)
				break
			}
			fmt.Println("unable to read from websocket:", err)
			continue
		}

		copy = c.handleStreamMessage(in, wsc)
		if copy {
			err = c.pushSocket.Send(in)
			if err != nil {
				fmt.Println("unable to send to push socket:", err)
			}
		}

	}


}

func (c *Control) userAllowedToPublish() bool {
	if c.currentStream.PublishUser != c.userId {
		return false
	}
	return true
}


// return true if message should be copied to push socket
func (c *Control) handleStreamMessage(msg []byte, conn *ws.Conn) bool {
	var err error
	var stream *Stream
	var subSocket mangos.Socket
	streamMessage := message.GetRootAsStreamMessage(msg, 0)
	streamId := string(streamMessage.StreamId())

	stream, err = GetStream(streamId)
	if err != nil {
		fmt.Println("unable to get stream:", err)
		c.loopback <- c.GetStreamErrorMessage(err)
	}

	fmt.Println("message lag:", GetCurrentTimeInMilli() - streamMessage.Timestamp())

	switch streamMessage.MessageType() {
	case message.MessageStreamBroadCast:
		// if user is allowed to broadcast on this stream
		if c.currentStream != stream && c.userAllowedToPublish() {
			c.currentStream = stream
			// setup push socket
			c.pushSocket, err = c.SetupPushSocket(stream)
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				c.loopback <- c.GetStreamErrorMessage(err)
			}

			subSocket, err = c.SetupSubSocket(stream)
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				c.loopback <- c.GetStreamErrorMessage(err)
			}
			// set publish true
			c.publish = true
			// close existing sub socket
			c.closeSubSocket <- true
			// start background write subscribe socket to web socket
			go c.CopyToWS(conn, subSocket)
		} else {
			fmt.Println("user not allowed to broadcase")
			c.loopback <- c.GetStreamErrorMessage(STREAM_NOT_ALLOWED)
		}
		//TODO what to do with duplicate broadcast messages ?
		return false
	case message.MessageStreamPause, message.MessageStreamStop, message.MessageStreamFrame:
		if !c.publish {
			c.loopback <- c.GetStreamErrorMessage(STREAM_NOT_ALLOWED)
			return false
		} else {
			return true
		}
	case message.MessageStreamSubscribe:
		if c.currentStream != stream {
			c.currentStream = stream
			// setup push socket
			c.pushSocket, err = c.SetupPushSocket(stream)
			if err != nil {
				fmt.Println("unable to setup push socket:", err)
				c.loopback <- c.GetStreamErrorMessage(err)
			}

			subSocket, err = c.SetupSubSocket(stream)
			if err != nil {
				fmt.Println("unable to setup sub socket:", err)
				c.loopback <- c.GetStreamErrorMessage(err)
			}

			c.publish = false
			// close existing sub socket
			c.closeSubSocket <- true

			go c.CopyToWS(conn, subSocket)
		}
		return false
	case message.MessageStreamComment:
		return true
	default:
		return false
	}
}

func (c *Control) GetStreamErrorMessage(err error) []byte {
	builder := fb.NewBuilder(0)
	message.StreamResponseStart(builder)

	switch err {
	case nil:
		message.StreamResponseAddStatus(builder, message.StatusOK)
	case STREAM_NOT_FOUND:
		message.StreamResponseAddStatus(builder, message.StatusNoStream)
	case STREAM_NOT_ALLOWED:
		message.StreamResponseAddStatus(builder, message.StatusNotAllowed)
	default:
		message.StreamResponseAddStatus(builder, message.StatusError)
	}

	message.StreamResponseEnd(builder)
	return builder.FinishedBytes()
}

