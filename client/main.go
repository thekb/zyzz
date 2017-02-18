package main

import (
	"github.com/gordonklaus/portaudio"
	"fmt"
	"gopkg.in/hraban/opus.v2"
	//"github.com/grd/ogg"
	"math/rand"
	"os"
	//"encoding/binary"
	//"golang.org/x/crypto/openpgp/packet"
	"github.com/grd/ogg"
	"encoding/binary"
	"github.com/thekb/zyzz/encode"
)

const (
	SAMPLE_RATE = 16000
	FRAMES_PER_BUFFER = 80 // opus legal frame size 2.5, 5, 10, 50 ms
	CHANNELS = 1
	CONTENT_TYPE = "audio/ogg; codecs=opus"
	HEADER_CONTENT_TYPE_OPTIONS = "X-Content-Type-Options"
	OPTION_NO_SNIFF = "nosniff"

)

func main() {

	var err error
	err = portaudio.Initialize()
	if err != nil {
		fmt.Println("error initializing port audio:", err)
		return
	}
	defer portaudio.Terminate()

	// 16 bit per sample
	input := make([]int16, FRAMES_PER_BUFFER)
	stream, err := portaudio.OpenDefaultStream(CHANNELS, 0, SAMPLE_RATE, FRAMES_PER_BUFFER, input)
	if err != nil {
		fmt.Println("unable to open default stream:", err)
		return
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		fmt.Println("unable to start stream:", err)
		return
	}
	var f *os.File
	f, err = os.Create("test.opus")
	if err != nil {
		fmt.Println("unable to create file:", err)
		return
	}
	defer f.Close()

	output := make([]byte, 1024)
	// init encoder
	var encoder opus.Encoder
	encoder.Init(SAMPLE_RATE, CHANNELS, opus.AppAudio)
	encoder.SetBitrateToAuto()
	//create ogg streamsate
	var streamState ogg.StreamState
	streamState.Init(rand.Int31())
	// create ogg packet
	var packet ogg.Packet
	// create ogg page
	var page ogg.Page

	header := encode.OpusHeader{
		Version: 1,
		Channels: 1,
		PreSkip: 0,
		InputSampleRate: SAMPLE_RATE,
		Gain: 1,
		ChannelMapping: 1,

	}
	packet.Packet = header.GetBytes()
	packet.BOS = true
	streamState.PacketIn(&packet)
	streamState.Flush(&page)
	binary.Write(f, binary.LittleEndian, page.Header)
	binary.Write(f, binary.LittleEndian, page.Body)
	var n int
	var currentPageSize int
	for k := 0; k < 1000; k++{
		err = stream.Read()
		if err != nil {
			fmt.Println("error reading stream:", err)
			break
		}
		n, err = encoder.Encode(input, output)
		if err != nil {
			fmt.Println("unable to encode:", err)
			break
		}
		currentPageSize += n
		packet.Packet = output[:n]
		packet.PacketNo += 1
		streamState.PacketIn(&packet)
		// if page contains more than 4k page out
		if currentPageSize > 1024 {
			streamState.PageOut(&page)
			currentPageSize = 0
			binary.Write(f, binary.BigEndian, page.Header)
			binary.Write(f, binary.BigEndian, page.Body)
		}


	}
	packet.EOS = true
	streamState.PacketIn(&packet)
	streamState.Flush(&page)
	binary.Write(f, binary.BigEndian, page.Header)
	binary.Write(f, binary.BigEndian, page.Body)
}
