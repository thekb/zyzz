package main

import (
	"github.com/gordonklaus/portaudio"
	"fmt"
	"gopkg.in/hraban/opus.v2"
	"github.com/grd/ogg"
	"math/rand"
	"os"
	"github.com/thekb/zyzz/encode"
	//"github.com/mccoyst/ogg"
)

const (
	SAMPLE_RATE = 16000
	FRAME_SIZE = 20 //opus legal frame size 2.5, 5, 10, 20 ms
	FRAMES_PER_BUFFER = (SAMPLE_RATE * FRAME_SIZE)/1000
	GRANULE_SAMPLES = (48000 * FRAME_SIZE)/1000
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
	var opusEncoder opus.Encoder
	opusEncoder.Init(SAMPLE_RATE, CHANNELS, opus.AppAudio)
	//opusEncoder.SetBitrate(32)
	//opusEncoder.SetBitrateToAuto()
	/*
	oggEncoder := ogg.NewEncoder(rand.Uint32(), f)
	*/
	var granulePosition int64

	opusHeader := encode.OpusHeader{
		Version: 1,
		Channels: 1,
		PreSkip: 0,
		InputSampleRate: SAMPLE_RATE,
		OutPutGain: 0,
		ChannelMap: 0,
	}

	opusCommentHeader := encode.OpusCommentHeader{
		VendorString: "thekbencoder",
		CommentList: []string{
			"NAME=stream",
			"ALBUM=album",
		},
	}
	/*
	oggEncoder.EncodeBOS(0, opusHeader.GetBytes())
	oggEncoder.Encode(0, opusCommentHeader.GetBytes())
	*/
	streamState := ogg.StreamState{}
	streamState.Init(rand.Int31())
	var packet ogg.Packet
	var page ogg.Page

	// write header
	packet.BOS = true
	packet.Packet = opusHeader.GetBytes()
	packet.GranulePos = granulePosition
	packet.PacketNo = 0
	fmt.Println(len(opusHeader.GetBytes()))
	fmt.Println(opusHeader.GetBytes())
	streamState.PacketIn(&packet)
	//streamState.PageOut(&page)
	streamState.Flush(&page)
	f.Write(page.Header)
	f.Write(page.Body)

	// write comment header
	packet.BOS = false
	packet.Packet = opusCommentHeader.GetBytes()
	packet.GranulePos = granulePosition
	packet.PacketNo = 1
	streamState.PacketIn(&packet)
	streamState.Flush(&page)
	f.Write(page.Header)
	f.Write(page.Body)
	// write header complete

	var n int
	granulePosition = GRANULE_SAMPLES
	for k := 0; k < 10;{
		err = stream.Read()
		if err != nil {
			fmt.Println("error reading stream:", err)
			break
		}
		n, err = opusEncoder.Encode(input, output)
		if err != nil {
			fmt.Println("unable to encode:", err)
			break
		}

		if n > 2 {
			packet.GranulePos = granulePosition
			packet.Packet = output[:n]
			packet.PacketNo += 1
			streamState.PacketIn(&packet)
			granulePosition += GRANULE_SAMPLES
			k++

		}
		//flush every 5 packets to file
		if k % 5 == 0 {
			streamState.Flush(&page)
			if len(page.Header) >0 && len(page.Body) > 0 {
				f.Write(page.Header)
				f.Write(page.Body)

			} else {
				fmt.Println("header or byd empty")
			}
		}
	}


	packet.EOS = true
	//packet.Packet = []byte{}
	//packet.GranulePos = 0
	streamState.PacketIn(&packet)
	streamState.Flush(&page)
	f.Write(page.Header)
	f.Write(page.Body)

	//oggEncoder.EncodeEOS()


}
