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
CHANNELS = 1

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
        try:
            chunk = stream.read(CHUNK)
            yield chunk
        except IOError as ioe:
            print "error %s" % ioe

url = "https://s1.zyzz.co/stream/publish/WGT1aN-CC/"

s = requests.session()
s.headers.update({'Content-Type': "audio/x-wav;codec=pcm"})
resp = s.post(url, data=get_chunks(STREAM))
