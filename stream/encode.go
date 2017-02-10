package stream

import (
	"github.com/winlinvip/go-fdkaac/fdkaac"
	"fmt"
)

const (
	CHANNELS = 2
	SAMPLE_RATE = 44100
	BITS_PER_SAMPLE = 16
	BIT_RATE = CHANNELS * SAMPLE_RATE * BITS_PER_SAMPLE

)

var aacEncoder *fdkaac.AacEncoder

func init() {
	aacEncoder = fdkaac.NewAacEncoder()
	fmt.Printf("frame size %v", aacEncoder.FrameSize())
	aacEncoder.InitLc(CHANNELS, SAMPLE_RATE, BIT_RATE)
	fmt.Printf("frame size %v", aacEncoder.FrameSize())
}




func EncodePCMToAAC(fragment []byte) ([]byte, error) {
	return aacEncoder.Encode(fragment)
}