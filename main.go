package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"github.com/thekb/zyzz/api"
	"github.com/thekb/zyzz/db"
	"fmt"
	"github.com/thekb/zyzz/stream"
)

func main() {
	d, err := db.GetDB()
	if err != nil {
		fmt.Println("unable to connect to db:", err)
		return
	}

	mySessions := sessions.New(sessions.Config{
		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "gosessionid"
		Cookie: "mysessionid",
		// base64 urlencoding,
		// if you have strange name cookie name enable this
		DecodeCookie: false,
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		Expires: 0,
		// the length of the sessionid's cookie's value
		CookieLength: 32,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it
		DisableSubdomainPersistence: false,
	})

	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(mySessions)
	// user api
	userApi := app.Party("/api/user")
	userApi.Handle("GET", "/", &api.CreateUser{api.Common{DB:d}})
	userApi.Handle("POST" , "/:id", &api.GetUser{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams", &api.GetUserStream{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams/current", &api.GetCurrentUserStream{api.Common{DB:d}})


	// event api
	eventApi := app.Party("/api/event")
	eventApi.Handle("POST", "/", &api.CreateEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/:id", api.GetEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/", api.GetEvents{api.Common{DB:d}})
	eventApi.Handle("GET", "/streams", &api.GetEventStreams{api.Common{DB:d}})

	// auth api
	authapi := app.Party("/auth")
	authapi.Get("/:provider", api.Authenticate)
	//authapi.HandleFunc("GET", "/:provider", api.Authenticate)
	authapi.Handle("GET", "/:provider/callback", &api.FacebookCallback{api.Common{DB:d}})


	//stream server api
	streamServerApi := app.Party("/api/streamserver")
	streamServerApi.Handle("POST", "/", &api.CreateStreamServer{api.Common{DB:d}})
	streamServerApi.Handle("GET", "/:id", &api.GetStreamServer{api.Common{DB:d}})

	// stream api
	streamApi := app.Party("/api/stream")
	streamApi.Handle("POST" ,"/", &api.CreateStream{api.Common{DB:d}})
	streamApi.Handle("GET" ,"/", &api.GetStreams{api.Common{DB:d}})

	streamParty := app.Party("/stream")
	streamParty.Handle("GET", "/ws/publish/:id", &stream.PublishStream{api.Common{DB:d}})
	streamParty.Handle("GET", "/ws/subscribe/:id", &stream.WebSocketSubscriber{api.Common{DB:d}})
	streamParty.Handle("GET", "/http/subscribe/:id", &stream.SubscribeStream{api.Common{DB:d}})



	app.Listen(":8000")

}
