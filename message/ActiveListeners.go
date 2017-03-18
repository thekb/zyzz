// automatically generated, do not modify

package message

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ActiveListeners struct {
	_tab flatbuffers.Table
}

func (rcv *ActiveListeners) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ActiveListeners) ActiveListeners() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func ActiveListenersStart(builder *flatbuffers.Builder) { builder.StartObject(1) }
func ActiveListenersAddActiveListeners(builder *flatbuffers.Builder, activeListeners int32) { builder.PrependInt32Slot(0, activeListeners, 0) }
func ActiveListenersEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
