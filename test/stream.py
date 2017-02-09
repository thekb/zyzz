import requests
import time

url = "http://localhost:8000/api/stream/S7AjZ0TBd"

file_path = "/home/thekb/work/src/github.com/thekb/thanos/output.aac"


def get_chunks(f):
    while True:
        time.sleep(0.1)
        data = f.read(1024)
        if not data:
            print "end of input stream"
            f.close()
            break
        yield data


f = open(file_path, 'rb')

s = requests.session()
s.headers.update({'Content-Type': "audio/aac"})
resp = s.put(url, data=get_chunks(f))
print resp