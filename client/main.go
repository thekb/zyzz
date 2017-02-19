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
	"encoding/binary"
	"github.com/mccoyst/ogg"
	"github.com/thekb/zyzz/encode"
	"golang.org/x/crypto/openpgp/packet"
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
	var opusEncoder opus.Encoder
	opusEncoder.Init(SAMPLE_RATE, CHANNELS, opus.AppAudio)
	opusEncoder.SetBitrateToAuto()

	oggEncoder := ogg.NewEncoder(rand.Int31(), f)
	oggEncoder.EncodeEOS()
	// add header


	var n int
	for k := 0; k < 1000; k++{
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

	}
}
