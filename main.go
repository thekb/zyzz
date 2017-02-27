package main

import (
	"gopkg.in/kataras/iris.v6/adaptors/sessions/sessiondb/redis/service"
	"gopkg.in/kataras/iris.v6/adaptors/sessions/sessiondb/redis"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"github.com/thekb/zyzz/stream"
	"github.com/thekb/zyzz/api"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/db"
	"fmt"
)



func myMiddleware(ctx *iris.Context){
	if _, err := ctx.Session().GetInt("id"); err != nil {
		ctx.Redirect("/")
		return // don't call original handler
	}
	ctx.Next()
}

func main() {
	d, err := db.GetDB()
	if err != nil {
		fmt.Println("unable to connect to db:", err)
		return
	}

	mySessions := sessions.New(sessions.Config{
		Cookie: "shortsess",
		DecodeCookie: false,
		Expires: 0,
		CookieLength: 32,
		DisableSubdomainPersistence: false,
	})
	config := service.DefaultConfig()
	mySessions.UseDatabase(redis.New(config))
	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(mySessions)

	// auth api
	authapi := app.Party("/auth")
	authapi.Get("/:provider", api.Authenticate)
	authapi.Handle("GET", "/:provider/callback", &api.FacebookCallback{api.Common{DB:d}})

	// user api
	userApi := app.Party("/api/user", myMiddleware)
	userApi.Handle("GET", "/", &api.CreateUser{api.Common{DB:d}})
	userApi.Handle("POST" , "/:id", &api.GetUser{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams", &api.GetUserStream{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams/current", &api.GetCurrentUserStream{api.Common{DB:d}})


	// event api
	eventApi := app.Party("/api/event", myMiddleware)
	eventApi.Handle("POST", "/", &api.CreateEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/:id", api.GetEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/", api.GetEvents{api.Common{DB:d}})
	eventApi.Handle("GET", "/streams", &api.GetEventStreams{api.Common{DB:d}})

	//stream server api
	streamServerApi := app.Party("/api/streamserver", myMiddleware)
	streamServerApi.Handle("POST", "/", &api.CreateStreamServer{api.Common{DB:d}})
	streamServerApi.Handle("GET", "/:id", &api.GetStreamServer{api.Common{DB:d}})

	// stream api
	streamApi := app.Party("/api/stream", myMiddleware)
	streamApi.Handle("POST" ,"/", &api.CreateStream{api.Common{DB:d}})
	streamApi.Handle("GET" ,"/", &api.GetStreams{api.Common{DB:d}})

	streamParty := app.Party("/stream", myMiddleware)
	streamParty.Handle("GET", "/ws/publish/:id", &stream.PublishStream{api.Common{DB:d}})
	streamParty.Handle("GET", "/ws/subscribe/:id", &stream.WebSocketSubscriber{api.Common{DB:d}})
	streamParty.Handle("GET", "/http/subscribe/:id", &stream.SubscribeStream{api.Common{DB:d}})



	app.Listen(":8000")

}
