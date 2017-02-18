package encode

import (
	"bytes"
	"encoding/binary"
)

type OpusHeader struct {
	Version uint8
	Channels uint8
	PreSkip int16
	InputSampleRate uint32
	Gain int16
	ChannelMapping uint8
}

// returns header bytes
func (oh *OpusHeader) GetBytes() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, []byte{'O', 'p', 'u', 's', 'H', 'e', 'a', 'd'})
	binary.Write(buffer, binary.BigEndian, oh.Version)
	binary.Write(buffer, binary.BigEndian, oh.Channels)
	binary.Write(buffer, binary.BigEndian, oh.PreSkip)
	binary.Write(buffer, binary.BigEndian, oh.InputSampleRate)
	binary.Write(buffer, binary.BigEndian, oh.Gain)
	binary.Write(buffer, binary.BigEndian, oh.ChannelMapping)
	return buffer.Bytes()
}

