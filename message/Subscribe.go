// automatically generated by the FlatBuffers compiler, do not modify

package message

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Subscribe struct {
	_tab flatbuffers.Table
}

func GetRootAsSubscribe(buf []byte, offset flatbuffers.UOffsetT) *Subscribe {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Subscribe{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Subscribe) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Subscribe) Table() flatbuffers.Table {
	return rcv._tab
}

func SubscribeStart(builder *flatbuffers.Builder) {
	builder.StartObject(0)
}
func SubscribeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
