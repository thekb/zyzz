// automatically generated by the FlatBuffers compiler, do not modify

package message

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SetupStreamResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsSetupStreamResponse(buf []byte, offset flatbuffers.UOffsetT) *SetupStreamResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SetupStreamResponse{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *SetupStreamResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SetupStreamResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SetupStreamResponse) StreamId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *SetupStreamResponse) Status() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 1
}

func (rcv *SetupStreamResponse) MutateStatus(n int8) bool {
	return rcv._tab.MutateInt8Slot(6, n)
}

func SetupStreamResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func SetupStreamResponseAddStreamId(builder *flatbuffers.Builder, streamId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(streamId), 0)
}
func SetupStreamResponseAddStatus(builder *flatbuffers.Builder, status int8) {
	builder.PrependInt8Slot(1, status, 1)
}
func SetupStreamResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
