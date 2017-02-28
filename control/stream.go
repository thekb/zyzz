package control

import (
	"github.com/thekb/zyzz/message"
	fb "github.com/google/flatbuffers/go"
	"fmt"
)

// when true tells the handler to setup subscriber for stream
func HandleStreamMessage(msg []byte) bool {
	streamMessage := message.GetRootAsStreamMessage(msg, 0)
	streamId := string(streamMessage.StreamId())
	messageTable := new(fb.Table)
	streamMessage.Message(messageTable)

	switch streamMessage.MessageType() {
	case message.MessageStreamControl:
		controlMessage := new(message.StreamControl)
		controlMessage.Init(messageTable.Bytes, messageTable.Pos)
		return handleStreamControl(streamId, controlMessage)
	case message.MessageStreamFrame:
		streamFrame := new(message.StreamFrame)
		streamFrame.Init(messageTable.Bytes, messageTable.Pos)
		return handleStreamFrame(streamId, streamFrame)
	case message.MessageStreamComment:
		streamComment := new(message.StreamComment)
		streamComment.Init(messageTable.Bytes, messageTable.Pos)
		return handleStreamComment(streamId, streamComment)
	default:
		fmt.Println("invalid message")
	}
	return false
}

func handleStreamControl(streamId string, controlMessage *message.StreamControl) bool {
	fmt.Println("reveived stream control message")
	switch controlMessage.StreamAction() {
	case message.StreamActionStart:
		fmt.Println("received stream start")
		//set stream state started
	case message.StreamActionPause:
		fmt.Println("received stream pause")
		//set stream state paused and broadcast to subscribers
	case message.StreamActionStop:
		fmt.Println("received stream stop")
		//set stream state stopped and broadcast to subscribers
	case message.StreamActionSubscribe:
		fmt.Println("received stream subscribe")
		// increment stream subscriber count
		return true
	}
	return false

}

func handleStreamFrame(streamId string, streamFrame *message.StreamFrame) bool {
	fmt.Println("reveived stream frame message")
	// publish stream bytes to socket
	streamFrame.FrameBytes()
	return false
}

func handleStreamComment(streamId string, streamComment *message.StreamComment) bool {
	return false
}