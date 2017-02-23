package packet

import (
	"errors"
	"bytes"
	"encoding/binary"
)

// packet format
// Zyzz --> magic string 4 -> 4 bytes
// timestamp --> int64 --> 8 bytes little endian --> unix time in nano seconds
// stream id --> byte array 12 bytes --> stream id --> stream xid
// user id --> byte array 12 bytes --> user xid
// sampling rate --> uint32 --> 4 bytes little endian --> input sampling rate
// channels --> uint8 --> 1 byte little endian --> input channels
// frame number --> uint64 --> 8 bytes little endian --> frame position reported by publisher
// input codec --> uint8 --> 1 byte little endian --> 0 - pcm/ 1 - opus/2 - aac
// frame size --> uint32 --> 8 bytes little endian --> frame size
// frame --> byte array


const (
	MAGIC_STRING = []byte{'Z', 'y', 'z', 'z'}
	INVALID_PACKET = errors.New("Invalid Packet")
)

type Packet struct {
	Timestamp int64
	StreamId string
	UserId string
	SamplingRate uint32
	Channels uint8
	FrameNumber uint64
	InputCodec uint8
	FrameSize uint32
	Frame []byte
}

func (p *Packet) decode(input []byte) error {
	// verify if input has magic string
	if !bytes.Equal(input[0:4], MAGIC_STRING) {
		return INVALID_PACKET
	}
	// timestamp offset 4
	p.Timestamp = int64(binary.LittleEndian.Uint64(input[4:12]))
	// stream id offset 12
	p.StreamId = string(input[12:24])
	// user id offset 24
	p.UserId = string(input[24:36])
	// sampling rate offset 36
	p.SamplingRate = binary.LittleEndian.Uint32(input[36:40])
	// channels offset 40
	p.Channels = uint8(input[40:41])
	// frame number 41
	p.FrameNumber = binary.LittleEndian.Uint64(input[41:49])
	// input code offset 49
	p.InputCodec = uint8(input[49:50])
	// frame size offset 50
	p.FrameSize = binary.LittleEndian.Uint32(input[50:54])
	// frame offset 54, read until frame size
	p.Frame = input[54:p.FrameSize]
	return nil
}

func GetPacket(input []byte) Packet {
	return nil
}