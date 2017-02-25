package main

import (
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/api"
	"github.com/thekb/zyzz/db"
	"fmt"
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
	userApi.Handle("POST", "/", api.CreateUser{api.Common{DB:d}})
	userApi.Handle("GET", "/:id", api.GetUser{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams", api.GetUserStream{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams/current", api.GetCurrentUserStream{api.Common{DB:d}})

	// event api
	eventApi := app.Party("/api/event")
	eventApi.Handle("POST", "/", api.CreateEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/:id", api.GetEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/", api.GetEvents{api.Common{DB:d}})
	eventApi.Handle("GET", "/streams", api.GetEventStreams{api.Common{DB:d}})

	// auth api
	authapi := app.Party("/auth")
	//authapi.Handle("GET", "/", api.SelectAuthenticator{api.Common{DB:d}})
	//authapi.Handle("GET", "/:provider", api.Authenticate)
	authapi.Get("/:provider", api.Authenticate)
	authapi.Handle("GET", "/:provider/callback", api.Callback{api.Common{DB:d}})
	/*
	// stream server api
	streamServerApi := app.Party("/api/streamserver")
	streamServerApi.Handle("POST", "/", api.CreateStreamServer{api.Common{DB:d}})
	streamServerApi.Handle("GET", "/:id", api.GetStreamServer{api.Common{DB:d}})

	// stream api
	streamApi := app.Party("/api/stream")
	streamApi.Handle("POST", "/", api.CreateStream{api.Common{DB:d}})
	streamApi.Handle("GET", "/:id", api.GetEvent{api.Common{DB:d}})
	streamApi.Handle("GET", "/", api.GetEvents{api.Common{DB:d}})

	// stream
	publishWS := websocket.New(websocket.Config{
		Endpoint: "/stream/publish/:id",
	})
	publishWS.OnConnection()

	subscribeWS := websocket.New(websocket.Config{
		Endpoint: "/stream/subscribe/:id",
	})
	app.Adapt(publishWS)
	app.Adapt(subscribeWS)
	*/
	app.Listen(":8000")

}
