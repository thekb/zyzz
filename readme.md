convert input to 64k mono aac output
1. ffmpeg -i input.file -ac 1 -c:a aac -b:a 64k output.aac

should recompile ffmpeg with libfdk_aac to support streaming

api server + multiple streaming servers

api server functions :

1. maintain list of available streams + metadata 
(number of concurrent users, when started etc...)
2. least connection load balancing between streaming servers
3. create stream 
4. 

streaming publish/subscribe will always go the same server

api should redirect  to appropriate streaming server

zzyz.co -> s1.zzyz.co, s2.zzyz.co, ...

db requirements fo
1. list of streams
2. stream metadata (Name, Publisher, active, when published, stream server)
3. 

```
apis
POST http://localhost:8000/api/user/

{"name":"test user", "description":"test user"}

{
    "error": "",
    "data": {
        "name": "test user",
        "description": "test user",
        "id": "I1TInx--C",
        "created_at": "2017-02-15T10:03:18Z",
        "published": 0,
        "subscribed": 0
    }
}

GET http://localhost:8000/api/user/I1TInx--C/

{
    "error": "",
    "data": {
        "name": "test user",
        "description": "test user",
        "id": "I1TInx--C",
        "created_at": "2017-02-15T10:03:18Z",
        "published": 0,
        "subscribed": 0
    }
}

POST http://localhost:8000/api/streamserver/

{"name": "stream server 1", "host_name": "s1.zyzz.co", "internal_ip":"172.31.27.248", "external_ip":"35.154.152.224"}


{
    "error": "",
    "data": {
        "id": "f-3vOxC-C",
        "name": "stream server 1",
        "host_name": "s1.zyzz.co",
        "internal_ip": "172.31.27.248",
        "external_ip": "35.154.152.224"
    }
}

POST http://localhost:8000/api/stream/

{"name":"stream 1", "description":"stream 1"}

{
    "error": "",
    "data": {
        "id": "OKkH7T--C",
        "name": "stream 1",
        "description": "stream 1",
        "started_at": "2017-02-15T10:17:35Z",
        "ended_at": "2017-02-15T10:17:35Z",
        "status": 0,
        "endpoint": "https://s1.zyzz.co/stream/OKkH7T--C/",
        "subscriber_count": 0,
        "creator_id": 1,
        "stream_server_id": 1,
        "transport_url": "ipc:///tmp/stream_OKkH7T--C.ipc"
    }
}

GET http://localhost:8000/api/stream/

{
    "error": "",
    "data": [{
        "id": "FbOoOTC--",
        "name": "stream 1",
        "description": "stream 1",
        "started_at": "2017-02-15T10:15:29Z",
        "ended_at": "2017-02-15T10:15:29Z",
        "status": 0,
        "endpoint": "https://s1.zyzz.co/stream/FbOoOTC--/",
        "subscriber_count": 0,
        "creator_id": 1,
        "stream_server_id": 1,
        "transport_url": "ipc:///tmp/stream_FbOoOTC--.ipc"
    }, {
        "id": "DgTK7TC--",
        "name": "stream 1",
        "description": "stream 1",
        "started_at": "2017-02-15T10:16:58Z",
        "ended_at": "2017-02-15T10:16:58Z",
        "status": 0,
        "endpoint": "https://s1.zyzz.co/stream/DgTK7TC--/",
        "subscriber_count": 0,
        "creator_id": 1,
        "stream_server_id": 1,
        "transport_url": "ipc:///tmp/stream_DgTK7TC--.ipc"
    }, {
        "id": "OKkH7T--C",
        "name": "stream 1",
        "description": "stream 1",
        "started_at": "2017-02-15T10:17:35Z",
        "ended_at": "2017-02-15T10:17:35Z",
        "status": 0,
        "endpoint": "https://s1.zyzz.co/stream/OKkH7T--C/",
        "subscriber_count": 0,
        "creator_id": 1,
        "stream_server_id": 1,
        "transport_url": "ipc:///tmp/stream_OKkH7T--C.ipc"
    }]
}

GET http://localhost:8000/api/stream/FbOoOTC--/

{
    "error": "",
    "data": {
        "id": "FbOoOTC--",
        "name": "stream 1",
        "description": "stream 1",
        "started_at": "2017-02-15T10:15:29Z",
        "ended_at": "2017-02-15T10:15:29Z",
        "status": 0,
        "endpoint": "https://s1.zyzz.co/stream/FbOoOTC--/",
        "subscriber_count": 0,
        "creator_id": 1,
        "stream_server_id": 1,
        "transport_url": "ipc:///tmp/stream_FbOoOTC--.ipc"
    }
}


```
required apis
1. facebook auth
2. create user after facebook auth
3. maitain session after login
4. user stream apis
    1. current published stream for user
    2. all streams published by user
5. create event
6. active streams on event

required features
1. offline record stream
2. chat on event stream
3. social media share




