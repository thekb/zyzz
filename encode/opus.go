package encode

import (
	"bytes"
	"encoding/binary"
	"github.com/grd/ogg"
	"io"
	"fmt"
)

const VENDOR_STRING  = "zyzz opus encoder 0.1"

type OpusOggStream struct {
	StreamId        int32
	Channels        uint8
	PreSkip         uint16
	InputSampleRate uint32
	OutPutGain      int16
	ChannelMap      uint8 //0 for mono or stereo
	VendorString    string
	CommentList     map[string]string
	streamState     ogg.StreamState
	granulePosition int64
	packet          ogg.Packet
	page            ogg.Page
}

func (oos *OpusOggStream) Start(writer io.Writer) {
	oos.streamState.Init(oos.StreamId)
	buffer := new(bytes.Buffer)
	// create header packet
	buffer.Write([]byte{
		'O', 'p', 'u', 's', 'H', 'e', 'a', 'd', //magic signature
		0x01, // version number

	})
	buffer.Write([]byte{oos.Channels}) // channel count
	binary.Write(buffer, binary.LittleEndian, oos.PreSkip) // pre skip
	binary.Write(buffer, binary.LittleEndian, oos.InputSampleRate) //input sample rate
	binary.Write(buffer, binary.LittleEndian, oos.OutPutGain) // output gain
	buffer.Write([]byte{oos.ChannelMap}) //channel mapping family
	//TODO implement channel mapping table
	//write header packet
	oos.packet.Packet = buffer.Bytes()
	oos.packet.BOS = true
	oos.granulePosition = 0
	oos.streamState.PacketIn(&oos.packet)
	oos.streamState.Flush(&oos.page)
	writer.Write(oos.page.Header)
	writer.Write(oos.page.Body)
	// create tags packet
	buffer.Reset()
	buffer.Write([]byte{'O', 'p', 'u', 's', 'T', 'a', 'g', 's'}) //magic signature
	// write vendor string length
	binary.Write(buffer, binary.LittleEndian, uint32(len(VENDOR_STRING)))
	//write vendor string
	buffer.Write([]byte(VENDOR_STRING))
	//write comment list length
	binary.Write(buffer, binary.LittleEndian, uint32(len(oos.CommentList)))
	for key, value := range oos.CommentList {
		comment := fmt.Sprintf("%s=s", key, value)
		binary.Write(buffer, binary.LittleEndian, uint32(len(comment)))
		buffer.Write([]byte(comment))
	}
	oos.packet.Packet = buffer.Bytes()
	oos.packet.BOS = false
	oos.granulePosition = 0
	oos.streamState.PacketIn(&oos.packet)
	oos.streamState.Flush(&oos.page)
	writer.Write(oos.page.Header)
	writer.Write(oos.page.Body)
}

func (oos *OpusOggStream) WritePacket(opusPacket []byte) {

}


type OpusHeader struct {
	Version         uint8
	Channels        uint8
	PreSkip         uint16
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
	binary.Write(buffer, binary.LittleEndian, oh.PreSkip)
	binary.Write(buffer, binary.LittleEndian, oh.InputSampleRate)
	binary.Write(buffer, binary.LittleEndian, oh.OutPutGain)
	buffer.Write([]byte{oh.ChannelMap})
	//buffer.Write([]byte{0x00})
	return buffer.Bytes()
}

type OpusCommentHeader struct {
	VendorString string
	CommentList  []string
}

func (och *OpusCommentHeader) GetBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte{'O', 'p', 'u', 's', 'T', 'a', 'g', 's'})
	// write vendor string length
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
