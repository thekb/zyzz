# automatically generated by the FlatBuffers compiler, do not modify

# namespace: message

import flatbuffers

class StreamStop(object):
    __slots__ = ['_tab']

    @classmethod
    def GetRootAsStreamStop(cls, buf, offset):
        n = flatbuffers.encode.Get(flatbuffers.packer.uoffset, buf, offset)
        x = StreamStop()
        x.Init(buf, n + offset)
        return x

    # StreamStop
    def Init(self, buf, pos):
        self._tab = flatbuffers.table.Table(buf, pos)

def StreamStopStart(builder): builder.StartObject(0)
def StreamStopEnd(builder): return builder.EndObject()
