package stream

import (
	"github.com/winlinvip/go-fdkaac/fdkaac"
)

const (
	CHUNK = 1024
	CHANNELS = 1
	SAMPLE_RATE = 44100
	BITS_PER_SAMPLE = 16
	BIT_RATE = CHANNELS * SAMPLE_RATE * BITS_PER_SAMPLE

)


func GetNewEncoder() *fdkaac.AacEncoder {
	aacEncoder := fdkaac.NewAacEncoder()
	aacEncoder.InitLc(CHANNELS, SAMPLE_RATE, BIT_RATE)
	return aacEncoder
}
