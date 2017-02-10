import pyaudio
import requests
import time


p = pyaudio.PyAudio()
# frames per buffer ?
CHUNK = 1024
# 16 bits per sample ?
FORMAT = pyaudio.paInt16
# 44.1k sampling rate ?
RATE = 44100
# number of channels
CHANNELS = 2

STREAM = p.open(
    format=FORMAT,
    channels=CHANNELS,
    rate=RATE,
    input=True,
    frames_per_buffer=CHUNK
)
print "initialized stream"


def get_chunks(stream):
    while True:
        #time.sleep(0.1)
        chunk = stream.read(CHUNK)
        #print len(chunk)

        yield chunk

url = "http://localhost:8000/api/stream/er18MPTdB"

s = requests.session()
s.headers.update({'Content-Type': "audio/x-wav;codec=pcm;rate=44100"})
resp = s.put(url, data=get_chunks(STREAM))
