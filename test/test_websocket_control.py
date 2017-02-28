from websocket import create_connection
import flatbuffers

import sys
sys.path.append(".")

from message import EventMessage
from message import StreamControl
from message.Message import Message
URL = "ws://localhost:8000/control"

ws = create_connection(URL)

print "connected"

builder = flatbuffers.Builder(0)
event_id_pos = builder.CreateString("event_id")
stream_id_pos = builder.CreateString("stream_id")
StreamControl.StreamControlStart(builder)
StreamControl.StreamControlAddStreamId(builder, stream_id_pos)
StreamControl.
stream_control_pos = StreamControl.StreamControlEnd(builder)

EventMessage.EventMessageStart(builder)
EventMessage.EventMessageAddEventId(builder, event_id_pos)
EventMessage.EventMessageAddMessageType(builder, Message.StreamControl)
EventMessage.EventMessageAddMessage(builder, stream_control_pos)
event_message_end = EventMessage.EventMessageEnd(builder)

builder.Finish(event_message_end)

ws.send_binary(builder.Output())

ws.close()


