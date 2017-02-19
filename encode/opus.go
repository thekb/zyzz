package encode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)


type OpusHeader struct {
	Version         uint8
	Channels        uint8
	PreSkip         int16
	InputSampleRate uint32
	OutPutGain      int16
	ChannelMap      uint8 //0 for mono or stereo
}

// returns header bytes
func (oh *OpusHeader) GetBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte{'O', 'p', 'u', 's', 'H', 'e', 'a', 'd'})
	buffer.Write([]byte{oh.Version})
	buffer.Write([]byte{oh.Channels})
	binary.Write(buffer, binary.LittleEndian, oh.Version)
	binary.Write(buffer, binary.LittleEndian, oh.InputSampleRate)
	binary.Write(buffer, binary.LittleEndian, oh.OutPutGain)
	buffer.Write([]byte{oh.ChannelMap})
	return buffer.Bytes()
}

type OpusCommentHeader struct {
	VendorString string
	CommentList []string
}

func (och *OpusCommentHeader) GetBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte{'O', 'p', 'u', 's', 'T', 'a', 'g', 's'})
	// write vendor string length
	fmt.Println(len(och.VendorString))
	binary.Write(buffer, binary.LittleEndian, uint32(len(och.VendorString)))
	//write vendor string
	buffer.Write([]byte(och.VendorString))
	//write comment list length
	binary.Write(buffer, binary.LittleEndian, uint32(len(och.CommentList)))
	for _, userComment := range och.CommentList {
		binary.Write(buffer, binary.LittleEndian, uint32(len(userComment)))
		buffer.Write([]byte(userComment))
	}
	return buffer.Bytes()
}
