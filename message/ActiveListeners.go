// automatically generated by the FlatBuffers compiler, do not modify

package message

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ActiveListeners struct {
	_tab flatbuffers.Table
}

func GetRootAsActiveListeners(buf []byte, offset flatbuffers.UOffsetT) *ActiveListeners {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ActiveListeners{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ActiveListeners) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ActiveListeners) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ActiveListeners) ActiveListeners() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ActiveListeners) MutateActiveListeners(n int32) bool {
	return rcv._tab.MutateInt32Slot(4, n)
}

func ActiveListenersStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ActiveListenersAddActiveListeners(builder *flatbuffers.Builder, activeListeners int32) {
	builder.PrependInt32Slot(0, activeListeners, 0)
}
func ActiveListenersEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
