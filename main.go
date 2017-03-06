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
	"github.com/thekb/zyzz/control"
	"gopkg.in/kataras/iris.v6/middleware/logger"
)



func sessionMiddleware(ctx *iris.Context){
	if _, err := ctx.Session().GetInt("id"); err != nil {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{})
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

	customLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
	})
	app.Use(customLogger)

	// auth api
	authapi := app.Party("/auth")
	authapi.Get("/:provider", api.Authenticate)
	authapi.Handle("GET", "/:provider/callback", &api.FacebookCallback{api.Common{DB:d}})

	// user api
	userApi := app.Party("/api/user", sessionMiddleware)
	userApi.Handle("GET", "/", &api.CreateUser{api.Common{DB:d}})
	userApi.Handle("POST" , "/:id", &api.GetUser{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams", &api.GetUserStream{api.Common{DB:d}})
	userApi.Handle("GET", "/:id/streams/current", &api.GetCurrentUserStream{api.Common{DB:d}})


	// event api
	eventApi := app.Party("/api/event")
	eventApi.Handle("POST", "/", &api.CreateEvent{api.Common{DB:d}})
	//eventApi.Handle("GET", "/:id", &api.GetEvent{api.Common{DB:d}})
	eventApi.Handle("GET", "/", &api.GetEvents{api.Common{DB:d}})
	eventApi.Handle("GET", "/:id/stream", &api.GetEventStreams{api.Common{DB:d}})
	eventApi.Handle("POST", "/:id/stream", &api.CreateStream{api.Common{DB:d}})

	//stream server api
	streamServerApi := app.Party("/api/streamserver")
	streamServerApi.Handle("POST", "/", &api.CreateStreamServer{api.Common{DB:d}})
	streamServerApi.Handle("GET", "/:id", &api.GetStreamServer{api.Common{DB:d}})


	streamParty := app.Party("/stream", sessionMiddleware)
	streamParty.Handle("GET", "/ws/publish/:id", &stream.PublishStream{api.Common{DB:d}})
	streamParty.Handle("GET", "/ws/subscribe/:id", &stream.WebSocketSubscriber{api.Common{DB:d}})
	streamParty.Handle("GET", "/http/subscribe/:id", &stream.SubscribeStream{api.Common{DB:d}})

	app.Handle("GET", "/control", &control.Control{})

	app.Listen(":8000")

}
