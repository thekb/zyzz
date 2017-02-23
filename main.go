package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"github.com/thekb/zyzz/api"
	"github.com/thekb/zyzz/db"
	"fmt"
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
)

func main() {
	d, err := db.GetDB()
	if err != nil {
		fmt.Println("unable to connect to db:", err)
		return
	}

	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	// user api
	userApi := app.Party("/api/user")
	userApi.Post("/", api.CreateUser{api.Common{DB:d}})
	userApi.Get("/:id", api.GetUser{api.Common{DB:d}})
	userApi.Get("/:id/streams", api.GetUserStream{api.Common{DB:d}})
	userApi.Get("/:id/streams/current", api.GetCurrentUserStream{api.Common{DB:d}})

	// event api
	eventApi := app.Party("/api/event")
	eventApi.Post("/", api.CreateEvent{api.Common{DB:d}})
	eventApi.Get("/:id", api.GetEvent{api.Common{DB:d}})
	eventApi.Get("/", api.GetEvents{api.Common{DB:d}})
	eventApi.Get("/streams", api.GetEventStreams{api.Common{DB:d}})

	// stream server api
	streamServerApi := app.Party("/api/streamserver")
	streamServerApi.Post("/", api.CreateStreamServer{api.Common{DB:d}})
	streamServerApi.Get("/:id", api.GetStreamServer{api.Common{DB:d}})

	// stream api
	streamApi := app.Party("/api/stream")
	streamApi.Post("/", api.CreateStream{api.Common{DB:d}})
	streamApi.Get("/:id", api.GetEvent{api.Common{DB:d}})
	streamApi.Get("/", api.GetEvents{api.Common{DB:d}})

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

	app.Listen(":8000")

}
