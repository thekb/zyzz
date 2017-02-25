package encode

import (
	"bytes"
	"encoding/binary"
	"github.com/grd/ogg"
	"io"
	"fmt"
)

const VENDOR_STRING = "zyzz opus encoder 0.1"

type OpusOggStream struct {
	StreamId        int32
	Channels        uint8
	PreSkip         uint16
	InputSampleRate uint32
	OutPutGain      int16
	ChannelMap      uint8 //0 for mono or stereo
	VendorString    string
	Comments        map[string]string
	FrameSize       float32
	streamState     ogg.StreamState
	granulePosition int64
	granuleSamples  int64
	packet          ogg.Packet
	page            ogg.Page
}

func (oos *OpusOggStream) Start(w io.Writer) {
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
	oos.granuleSamples = int64(48000 * oos.FrameSize /1000)
	oos.streamState.PacketIn(&oos.packet)
	oos.streamState.Flush(&oos.page)
	w.Write(oos.page.Header)
	w.Write(oos.page.Body)
	// create tags packet
	buffer.Reset()
	buffer.Write([]byte{'O', 'p', 'u', 's', 'T', 'a', 'g', 's'}) //magic signature
	// write vendor string length
	binary.Write(buffer, binary.LittleEndian, uint32(len(VENDOR_STRING)))
	//write vendor string
	buffer.Write([]byte(VENDOR_STRING))
	//write comment list length
	binary.Write(buffer, binary.LittleEndian, uint32(len(oos.Comments)))
	for key, value := range oos.Comments {
		comment := fmt.Sprintf("%s=s", key, value)
		binary.Write(buffer, binary.LittleEndian, uint32(len(comment)))
		buffer.Write([]byte(comment))
	}
	oos.packet.Packet = buffer.Bytes()
	oos.packet.BOS = false
	oos.granulePosition = 0
	oos.streamState.PacketIn(&oos.packet)
	oos.streamState.Flush(&oos.page)
	w.Write(oos.page.Header)
	w.Write(oos.page.Body)
}

// returns false when frame is flushed to writer
func (oos *OpusOggStream) WritePacket(opusPacket []byte, w io.Writer) bool {
	oos.granulePosition += oos.granuleSamples
	oos.packet.GranulePos = oos.granulePosition
	oos.packet.Packet = opusPacket
	oos.packet.PacketNo += 1
	oos.streamState.PacketIn(&oos.packet)
	oos.streamState.Flush(&oos.page)
	w.Write(oos.page.Header)
	w.Write(oos.page.Body)
	return false
}
