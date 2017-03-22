// automatically generated by the FlatBuffers compiler, do not modify

package message;

import java.nio.*;
import java.lang.*;
import java.util.*;
import com.google.flatbuffers.*;

@SuppressWarnings("unused")
public final class Response extends Table {
  public static Response getRootAsResponse(ByteBuffer _bb) { return getRootAsResponse(_bb, new Response()); }
  public static Response getRootAsResponse(ByteBuffer _bb, Response obj) { _bb.order(ByteOrder.LITTLE_ENDIAN); return (obj.__assign(_bb.getInt(_bb.position()) + _bb.position(), _bb)); }
  public void __init(int _i, ByteBuffer _bb) { bb_pos = _i; bb = _bb; }
  public Response __assign(int _i, ByteBuffer _bb) { __init(_i, _bb); return this; }

  public byte status() { int o = __offset(4); return o != 0 ? bb.get(o + bb_pos) : 1; }

  public static int createResponse(FlatBufferBuilder builder,
      byte status) {
    builder.startObject(1);
    Response.addStatus(builder, status);
    return Response.endResponse(builder);
  }

  public static void startResponse(FlatBufferBuilder builder) { builder.startObject(1); }
  public static void addStatus(FlatBufferBuilder builder, byte status) { builder.addByte(0, status, 1); }
  public static int endResponse(FlatBufferBuilder builder) {
    int o = builder.endObject();
    return o;
  }
}

