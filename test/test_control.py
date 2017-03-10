from websocket import create_connection
import websocket
import flatbuffers

import sys
sys.path.append(".")


URL = "ws://localhost:8000/control"

def on_open(ws):
    print "opened"

def on_close(ws):
    print "closed"

def on_message(ws, message):
    print message

def on_error(ws, error):
    print error

ws = websocket.WebSocketApp(URL,
                            on_open=on_open,
                            on_close=on_close,
                            on_message=on_message,
                            on_error=on_error,
                            header=["X-User-Id: 1", ]
                            )
ws.run_forever()
print "connected"



