# automatically generated by the FlatBuffers compiler, do not modify

# namespace: message

import flatbuffers

class StreamControl(object):
    __slots__ = ['_tab']

    @classmethod
    def GetRootAsStreamControl(cls, buf, offset):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, offset)
        x = StreamControl()
        x.Init(buf, n + offset)
        return x

    # StreamControl
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

    # StreamControl
    def SampleRate(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(4))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint32Flags, o + self._tab.Pos)
        return 0

    # StreamControl
    def Channels(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(6))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint8Flags, o + self._tab.Pos)
        return 0

    # StreamControl
    def FrameSize(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(8))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Uint8Flags, o + self._tab.Pos)
        return 0

    # StreamControl
    def StreamAction(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(10))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Int8Flags, o + self._tab.Pos)
        return 1

    # StreamControl
    def Encoding(self):
        o = flatbuffers.number_types.UOffsetTFlags.py_type(self._tab.Offset(12))
        if o != 0:
            return self._tab.Get(flatbuffers.number_types.Int8Flags, o + self._tab.Pos)
        return 1

def StreamControlStart(builder): builder.StartObject(5)
def StreamControlAddSampleRate(builder, sampleRate): builder.PrependUint32Slot(0, sampleRate, 0)
def StreamControlAddChannels(builder, channels): builder.PrependUint8Slot(1, channels, 0)
def StreamControlAddFrameSize(builder, frameSize): builder.PrependUint8Slot(2, frameSize, 0)
def StreamControlAddStreamAction(builder, streamAction): builder.PrependInt8Slot(3, streamAction, 1)
def StreamControlAddEncoding(builder, encoding): builder.PrependInt8Slot(4, encoding, 1)
def StreamControlEnd(builder): return builder.EndObject()
