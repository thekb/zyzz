// automatically generated by the FlatBuffers compiler, do not modify

package message

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type BroadCast struct {
	_tab flatbuffers.Table
}

func GetRootAsBroadCast(buf []byte, offset flatbuffers.UOffsetT) *BroadCast {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &BroadCast{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *BroadCast) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *BroadCast) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *BroadCast) Encoding() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 1
}

func (rcv *BroadCast) MutateEncoding(n int8) bool {
	return rcv._tab.MutateInt8Slot(4, n)
}

func BroadCastStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func BroadCastAddEncoding(builder *flatbuffers.Builder, encoding int8) {
	builder.PrependInt8Slot(0, encoding, 1)
}
func BroadCastEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
