package main

import (
	ws "github.com/gorilla/websocket"
	m "github.com/thekb/zyzz/message"
	"net/url"
	"fmt"
	pa "github.com/gordonklaus/portaudio"
	 "github.com/thekb/opus"
	fb "github.com/google/flatbuffers/go"
	"time"
	"net/http"
)

const (
	INPUT_CHANNELS = 1
	OUTPUT_CHANNELS = 0
	SAMPLE_RATE = 24000
	FRAMES_SIZE = 20 // in milliseconds
	FRAMES_PER_BUFFER = int(SAMPLE_RATE * FRAMES_SIZE/1000)
)


func main() {
	var err error
	var wsRead []byte
	streamId := "nRSFYNx--"
	eventId := "pfYX3Z1C-"
	b := fb.NewBuilder(1024)
	pa.Initialize()
	defer pa.Terminate()


	streamBuffer := make([]int16, FRAMES_PER_BUFFER)
	stream, err := pa.OpenDefaultStream(
		INPUT_CHANNELS, OUTPUT_CHANNELS,
		SAMPLE_RATE, FRAMES_PER_BUFFER, streamBuffer)
	defer stream.Close()


	url := url.URL{Scheme:"ws", Host: "localhost:8000", Path:"/control", }

	var header = make(http.Header)
	header.Set("X-User-Id", "1")

	c, _, err := ws.DefaultDialer.Dial(url.String(), header)

	if err != nil {
		fmt.Println("unable to connect to websocket:", err)
		return
	}

	// write broadcast message
	c.WriteMessage(ws.BinaryMessage, GetStreamBroadCastMessage(b, streamId, eventId))
	_, wsRead, err = c.ReadMessage()
	if err != nil {
		fmt.Println("unable to read from websocket:", err)
		return
	}

	streamMessage := m.GetRootAsStreamMessage(wsRead, 0)
	fmt.Println("stream mesage type:", streamMessage.MessageType())
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
	go func(c *ws.Conn){
		for {
			_, out, err := c.ReadMessage()
			if err != nil {
				fmt.Println("unable to read message from websocket:", err)
				break
			}
			message := m.GetRootAsStreamMessage(out, 0)
			table := new(fb.Table)
			fmt.Println(string(message.StreamId()))
			fmt.Println(message.MessageType())
			if message.Message(table) {
				if message.MessageType() == m.MessageComment {
					comment := new(m.Comment)
					comment.Init(table.Bytes, table.Pos)
					fmt.Println(string(comment.Text()))
					fmt.Println(string(comment.UserName()))
				}
			}
		}

	}(c)
	// loop will run for 1 second and break
	// each stream frame is followed by stream comment

	var n int
	var encoder *opus.Encoder
	var encoderBuffer = make([]byte, 1000)
	encoder, err = opus.NewEncoder(SAMPLE_RATE, INPUT_CHANNELS, opus.AppAudio)

	fmt.Println("sending stream")
	timeout := time.After(time.Second * 120)
	ticker := time.Tick(time.Millisecond * FRAMES_SIZE)

	L:
	for {
		select {
		case <- timeout:
			fmt.Println("after timeout")
			break L
		case <- ticker:
			fmt.Println("reading stream")
			stream.Read()
			n, err = encoder.Encode(streamBuffer, encoderBuffer)
			if err != nil {
				fmt.Println("unable to encode stream:", err)
				continue
			}
			if n > 2 {
				c.WriteMessage(ws.BinaryMessage, GetStreamFrameMessage(b, encoderBuffer[:n], streamId, eventId))
				c.WriteMessage(ws.BinaryMessage, GetCommentMessage(b, streamId, eventId, "username", "comment"))
			}

		}

	}

	// send stream pause
	c.WriteMessage(ws.BinaryMessage, GetStreamPauseMessage(b, streamId, eventId))
	// send comment
	c.WriteMessage(ws.BinaryMessage, GetCommentMessage(b, streamId, eventId, "username", "comment"))
	// send stream stop
	c.WriteMessage(ws.BinaryMessage, GetStreamStopMessage(b, streamId, eventId))
	// send comment
	c.WriteMessage(ws.BinaryMessage, GetCommentMessage(b, streamId, eventId, "username", "comment"))

}

func GetStreamPauseMessage(b *fb.Builder, streamId, eventId string) []byte {
	b.Reset()

	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)

	m.PauseStart(b)
	streamPauseOffset := m.PauseEnd(b)

	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessagePause)
	m.StreamMessageAddMessage(b, streamPauseOffset)
	m.StreamMessageAddTimestamp(b, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(b)

	b.Finish(streamMessageOffset)
	return b.FinishedBytes()
}

func GetStreamStopMessage(b *fb.Builder, streamId, eventId string) []byte {
	b.Reset()

	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)

	m.StopStart(b)
	streamStopOffset := m.StopEnd(b)

	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessageStop)
	m.StreamMessageAddMessage(b, streamStopOffset)
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

func GetStreamBroadCastMessage(b *fb.Builder, streamId, eventId string) []byte {
	b.Reset()
	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)

	m.BroadCastStart(b)
	m.BroadCastAddEncoding(b, m.InputEncodingOpus)
	streamBroadcastOffset := m.BroadCastEnd(b)

	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessageBroadCast)
	m.StreamMessageAddMessage(b, streamBroadcastOffset)
	m.StreamMessageAddTimestamp(b, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(b)

	b.Finish(streamMessageOffset)
	return b.FinishedBytes()
}

func GetStreamFrameMessage(b *fb.Builder, input []byte, streamId, eventId string) []byte {
	b.Reset()

	frameLength := len(input)

	m.FrameStartFrameVector(b, frameLength)
	// iterate in reverse order
	for i := frameLength - 1; i >=0; i-- {
		b.PrependByte(input[i])
	}
	frameOffset := b.EndVector(frameLength)

	streamIdOffset := b.CreateString(streamId)
	eventIdOffset := b.CreateString(eventId)

	m.FrameStart(b)
	m.FrameAddFrameSize(b, byte(FRAMES_SIZE))
	m.FrameAddSampleRate(b, uint32(SAMPLE_RATE))
	m.FrameAddChannels(b, byte(INPUT_CHANNELS))
	m.FrameAddFrame(b, frameOffset)
	streamFrameOffset := m.FrameEnd(b)
	m.StreamMessageStart(b)
	m.StreamMessageAddEventId(b, eventIdOffset)
	m.StreamMessageAddStreamId(b, streamIdOffset)
	m.StreamMessageAddMessageType(b, m.MessageFrame)
	m.StreamMessageAddMessage(b, streamFrameOffset)
	m.StreamMessageAddTimestamp(b, GetTimeInMillis())
	streamMessageOffset := m.StreamMessageEnd(b)
	b.Finish(streamMessageOffset)
	return b.FinishedBytes()
}

func GetTimeInMillis() int64 {
	 return time.Now().UnixNano() / int64(time.Millisecond)
}