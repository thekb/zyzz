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

{"name":"stream 1", "description":"stream 1", "event_id":1}

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
        "transport_url": "ipc:///tmp/stream_OKkH7T--C.ipc",
        "event_id": 1
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

POST http://localhost:8000/api/event/

{"name":"event 1", "description":"event 1", "running_now":0,"matchid":1027319, "matchurl":"http://synd.cricbuzz.com/j2me/1.0/match/2015/2015_WCLC/NEP_KEN_MAR11/"}


{
    "error": "",
    "data": {
        "name": "event 1",
        "description": "event 1",
        "id": "4UhQn5TCC",
        "created_at": "2017-03-11T21:09:22Z",
        "starttime": "0001-01-01T00:00:00Z",
        "endtime": "0001-01-01T00:00:00Z",
        "running_now": 0,
        "matchid": 1027319,
        "matchurl": "http://synd.cricbuzz.com/j2me/1.0/match/2015/2015_WCLC/NEP_KEN_MAR11/"
    }
}

GET http://localhost:8000/api/event/

{
    "error": "",
    "data": [{
        "name": "stream 1",
        "description": "event 1",
        "id": "wYYkFN1-C",
        "created_at": "2017-03-01T19:53:56Z",
        "starttime": "0001-01-01T00:00:00Z",
        "endtime": "0001-01-01T00:00:00Z",
        "running_now": 1
    }]
}
Cricket Store
GET http://localhost:8000/api/cricbuzz/4UhQn5TCC
{
    "error": "",
    "data": {
        "MatchInfo": {
            "Type": "ODI",
            "Srs": "ICC World Cricket League Championship, 2015-17",
            "MatchDesc": "NEP vs Ken",
            "MatchNumber": "33rd Match",
            "HostCity": "Kirtipur",
            "HostCountry": "Nepal",
            "Ground": "Tribhuvan University International Cricket Ground",
            "DataPath": "http://synd.cricbuzz.com/j2me/1.0/match/2015/2015_WCLC/NEP_KEN_MAR11/",
            "InngCnt": "",
            "MatchState": {
                "MatchState": "complete",
                "Status": "Ken won by 5 wickets (D/L method)",
                "TossWon": "Ken",
                "Decision": "Fielding",
                "AddnStatus": "",
                "SplStatus": ""
            },
            "Team": [{
                "Name": "NEP",
                "SName": "NEP",
                "Flag": "0"
            }, {
                "Name": "Ken",
                "SName": "KEN",
                "Flag": "1"
            }],
            "Schedule": {
                "StartTime": "03:45",
                "EndDate": "Mar 11 2017"
            },
            "Score": {
                "InningsDetail": {
                    "noOfOvers": "50",
                    "RequiredRunRate": "0",
                    "CurrentRunRate": "4.06",
                    "CurrentPartnership": ""
                },
                "BattingTeam": {
                    "SName": "KEN",
                    "Innings": [{
                        "Description": "Inns",
                        "Runs": "98",
                        "Declared": "0",
                        "FollowOn": "0",
                        "Overs": "24.1",
                        "Wickets": "5"
                    }]
                },
                "BowlingTeam": {
                    "SName": "NEP",
                    "Innings": [{
                        "Description": "Inns",
                        "Runs": "112",
                        "Declared": "0",
                        "FollowOn": "0",
                        "Overs": "36",
                        "Wickets": "8"
                    }]
                },
                "Batsmen": {
                    "SName": "",
                    "Runs": "12",
                    "Balls": "19",
                    "Fours": "1",
                    "Sixes": "0"
                },
                "Bowler": {
                    "SNames": "",
                    "Runs": "8",
                    "Wickets": "",
                    "Overs": "",
                    "Maidens": ""
                }
            }
        }
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

Starting up with postgres:
Install postgres and start postgres.
Run these commands once postgres is up and running.

psql
create database zyzz;
\q
psql zyzz
create user postgres with password 'melcow';

