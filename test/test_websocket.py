from websocket import create_connection
import pyaudio
import time
import struct
import array
import binascii

p = pyaudio.PyAudio()
# 16 bits per sample ?
FORMAT = pyaudio.paInt16
# 44.1k sampling rate ?
RATE = 24000
# number of channels
CHANNELS = 1
FRAME_SIZE = 60  # in milliseconds
# frames per buffer ?
CHUNK = int(RATE * FRAME_SIZE/1000)
print CHUNK
STREAM = p.open(
    format=FORMAT,
    channels=CHANNELS,
    rate=RATE,
    input=True,
    frames_per_buffer=CHUNK
)
print "initialized stream"

URL = "ws://localhost:8000/stream/uPeDTmCC-/publish"
i = 0


ws = create_connection(URL)
while True:
    chunk = STREAM.read(CHUNK)
    ws.send_binary(chunk)
    i += 1
