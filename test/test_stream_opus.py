from opuslib.api import encoder, constants, ctl
from opuslib.exceptions import OpusError
from websocket import create_connection
import pyaudio


p = pyaudio.PyAudio()
# 16 bits per sample ?
FORMAT = pyaudio.paInt16
# 44.1k sampling rate ?
RATE = 24000
# number of channels
CHANNELS = 1
FRAME_SIZE = 60  # in milliseconds
# frames per buffer ?
CHUNK_SIZE = int(RATE * FRAME_SIZE / 1000)
STREAM = p.open(
    format=FORMAT,
    channels=CHANNELS,
    rate=RATE,
    input=True,
    frames_per_buffer=CHUNK_SIZE
)
print "initialized stream"
enc = encoder.create(RATE, CHANNELS, constants.APPLICATION_AUDIO)

URL = "ws://localhost:8000/stream/ws/opus/publish/DVF_W31C-"

ws = create_connection(URL)
while True:
    chunk = STREAM.read(CHUNK_SIZE)
    encoded = encoder.encode(enc, chunk, CHUNK_SIZE, CHUNK_SIZE)
    print len(encoded)
    ws.send_binary(encoded)






