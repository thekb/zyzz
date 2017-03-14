package main

import (
	ws "github.com/gorilla/websocket"
	m "github.com/thekb/zyzz/message"
	fb "github.com/google/flatbuffers/go"
	"net/url"
	"fmt"
	"time"
)

func main () {
	var read []byte
	var c *ws.Conn
	var err error
	url := url.URL{Scheme:"ws", Host: "localhost:8000", Path:"/control", }

	c, _, err = ws.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println("unable to connect to websocket:", err)
		return
	}
	b := fb.NewBuilder(0)
	streamId := "nRSFYNx--"
	eventId := "pfYX3Z1C-"
	userName := "subscriber1"
	comment := "subscriber comment"

	// send subscribe message
	c.WriteMessage(ws.BinaryMessage, GetSubscribeMessage(b, eventId, streamId))
	// should receive ok
	_, read, err = c.ReadMessage()
	if err != nil {
		fmt.Println("unable to read message from websocket:", err)
		return
	}

	streamMessage := m.GetRootAsStreamMessage(read, 0)
	table := new(fb.Table)
	if streamMessage.Message(table) {
		if streamMessage.MessageType() == m.MessageResponse {
			response := new(m.Response)
			response.Init(table.Bytes, table.Pos)
			if response.Status() != m.ResponseStatusOK {
				fmt.Println("status not ok,", response.Status())
				c.Close()
				return
			}
		}
	}
	// read messages in background
	go func (c *ws.Conn){
		for {
			_, out, err := c.ReadMessage()
			if err != nil {
				fmt.Println("unable to read message from websocket:", err)
				break
			}
			fmt.Println(out)
			message := m.GetRootAsStreamMessage(out, 0)
			table := new(fb.Table)
			fmt.Println(message.MessageType())
			if message.Message(table) {
				switch message.MessageType() {
				case m.MessageComment:
					fmt.Println("received comment message")
					comment := new(m.Comment)
					comment.Init(table.Bytes, table.Pos)
					fmt.Println(string(comment.Text()))
					fmt.Println(string(comment.UserName()))
				case m.MessageFrame:
					fmt.Println("received frame message")
				case m.MessagePause:
					fmt.Println("received pause message")
				case m.MessageStop:
					fmt.Println("received stop message")
				case m.MessageStatus:
					fmt.Println("received status message")
					status := new(m.Status)
					status.Init(table.Bytes, table.Pos)
					fmt.Println("status:", status.Status())
				}

			}
		}
	}(c)

	timeout := time.After(time.Second * 2)
	ticker := time.Tick(time.Millisecond * 500)
	for {
		select {
		// un subscribe after 2 seconds
		case <- timeout:
			unSubscribeMsg := GetUnSubscribeMessage(b, eventId, streamId)
			c.WriteMessage(ws.BinaryMessage, unSubscribeMsg)
			fmt.Println("sent unsubscribe message")
		// send comment every 0.5 seconds
		case <- ticker:
			commentMsg := GetCommentMessage(b, streamId, eventId, userName, comment)
			c.WriteMessage(ws.BinaryMessage, commentMsg)
			fmt.Println("sent comment message")
			fmt.Println("tick")
		}
	}



}

func GetSubscribeMessage(b *fb.Builder, eventId, streamId string) []byte {
	b.Reset()

	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)

	m.SubscribeStart(b)
	subscribeOffset := m.SubscribeEnd(b)

	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessageSubscribe)
	m.StreamMessageAddMessage(b, subscribeOffset)
	m.StreamMessageAddTimestamp(b, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(b)

	b.Finish(streamMessageOffset)
	return b.FinishedBytes()
}

func GetUnSubscribeMessage(b *fb.Builder, eventId, streamId string) []byte {
	b.Reset()

	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)

	m.UnSubscribeStart(b)
	unSubscribeOffset := m.UnSubscribeEnd(b)

	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessageUnSubscribe)
	m.StreamMessageAddMessage(b, unSubscribeOffset)
	m.StreamMessageAddTimestamp(b, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(b)

	b.Finish(streamMessageOffset)
	return b.FinishedBytes()
}

func GetCommentMessage(b *fb.Builder, streamId, eventId, userName, text string) []byte {
	b.Reset()

	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)
	userNameOffset := b.CreateString(userName)
	textOffset := b.CreateString(text)

	m.CommentStart(b)
	m.CommentAddUserName(b, userNameOffset)
	m.CommentAddText(b, textOffset)
	streamCommentOffset := m.CommentEnd(b)

	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessageComment)
	m.StreamMessageAddMessage(b, streamCommentOffset)
	m.StreamMessageAddTimestamp(b, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(b)

	b.Finish(streamMessageOffset)
	return b.FinishedBytes()
}


func GetTimeInMillis() int64 {
	 return time.Now().UnixNano() / int64(time.Millisecond)
}
