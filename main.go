package main

import (
	rservice "gopkg.in/kataras/iris.v6/adaptors/sessions/sessiondb/redis/service"
	sredis "gopkg.in/kataras/iris.v6/adaptors/sessions/sessiondb/redis"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"gopkg.in/kataras/iris.v6/middleware/pprof"
	"gopkg.in/redis.v5"
	"github.com/thekb/zyzz/stream"
	"github.com/thekb/zyzz/api"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/db"
	"fmt"
	"github.com/thekb/zyzz/control"
	"gopkg.in/kataras/iris.v6/middleware/logger"
	"flag"
	"crypto/tls"
)

const (
	REDIS_NETWORK = "tcp"
	REDIS_ADDRESS = "127.0.0.1:6379"
	SESSION_PREFIX = "zyzz"
	CERT_PATH = "/etc/letsencrypt/live/shortpitch.live/fullchain.pem"
	KEY_PATH = "/etc/letsencrypt/live/shortpitch.live/privkey.pem"
)

func sessionMiddleware(ctx *iris.Context) {
	if _, err := ctx.Session().GetInt("id"); err != nil {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{})
		return // don't call original handler
	}
	ctx.Next()
}

func main() {
	// run in prod mode
	prod := flag.Bool("prod", false, "run in prod mode")
	// enable profiling
	debug := flag.Bool("debug", false, "enable profiling")

	flag.Parse()

	// get db instance
	d, err := db.GetDB()
	if err != nil {
		fmt.Println("unable to connect to db:", err)
		return
	}
	// initialize redis pool
	r := redis.NewClient(&redis.Options{
		Network: REDIS_NETWORK,
		Addr: REDIS_ADDRESS,
		PoolSize: 100,
	})
	_, err = r.Ping().Result()
	if err != nil {
		fmt.Println("unable to connect to redis:", err)
		return
	}
	// start event tickers in background
	api.StartEventTickers(d, r)

	// setup sessions
	session := sessions.New(sessions.Config{
		Cookie: "sid",
		DecodeCookie: false,
		Expires: 0,
		CookieLength: 32,
		DisableSubdomainPersistence: false,
	})
	session.UseDatabase(sredis.New(rservice.Config{
		Network: REDIS_NETWORK,
		Addr: REDIS_ADDRESS,
	}))

	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())
	app.Adapt(session)

	/*
	app.StaticWeb("/css","/Users/abalusu/Projects/shortpitch/dist/css")
	app.StaticWeb("/js","/Users/abalusu/Projects/shortpitch/dist/js")
	app.StaticWeb("/assets","/Users/abalusu/Projects/shortpitch/dist/assets")
	app.StaticWeb("/vendor","/Users/abalusu/Projects/shortpitch/dist/vendor")
	app.StaticWeb("/jspm_packages","/Users/abalusu/Projects/shortpitch/dist/jspm_packages")
	app.StaticWeb("/home","/Users/abalusu/Projects/shortpitch/dist/")
	app.HandleFunc("GET", "/", func(ctx *iris.Context) { ctx.Redirect("/home", 302) })
	*/
	app.StaticWeb("/css","/opt/shortpitch/UI/dist/css")
	app.StaticWeb("/js","/opt/shortpitch/UI/dist/js")
	app.StaticWeb("/assets","/opt/shortpitch/UI/dist/assets")
	app.StaticWeb("/vendor","/opt/shortpitch/UI/dist/vendor")
	app.StaticWeb("/jspm_packages","/opt/shortpitch/UI/dist/jspm_packages")
	app.StaticWeb("/home","/opt/shortpitch/UI/dist/")
	app.HandleFunc("GET", "/", func(ctx *iris.Context) { ctx.Redirect("/home", 302) })

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
	authapi.Handle("GET", "/:provider/callback", &api.FacebookCallback{api.Common{DB:d, R:r}})

	// user api
	userApi := app.Party("/api/user", sessionMiddleware)
	userApi.Handle("GET", "/", &api.CreateUser{api.Common{DB:d, R:r}})
	userApi.Handle("POST", "/:id", &api.GetUser{api.Common{DB:d, R:r}})
	userApi.Handle("GET", "/:id/streams", &api.GetUserStream{api.Common{DB:d, R:r}})
	userApi.Handle("GET", "/:id/streams/current", &api.GetCurrentUserStream{api.Common{DB:d, R:r}})


	// event api no auth
	eventApiNoAuth := app.Party("/api/event")
	eventApiNoAuth.Handle("GET", "/", &api.GetEvents{api.Common{DB:d, R:r}})
	eventApiNoAuth.Handle("GET", "/:id/stream", &api.GetEventStreams{api.Common{DB:d, R:r}})

	// event api auth
	eventApi := app.Party("/api/event", sessionMiddleware)
	eventApi.Handle("POST", "/", &api.CreateEvent{api.Common{DB:d, R:r}})
	eventApi.Handle("PUT", "/:id", &api.UpdateEvent{api.Common{DB:d, R:r}})
	eventApi.Handle("POST", "/:id/stream", &api.CreateStream{api.Common{DB:d, R:r}})

	//stream server api
	streamServerApi := app.Party("/api/streamserver")
	streamServerApi.Handle("POST", "/", &api.CreateStreamServer{api.Common{DB:d, R:r}})
	streamServerApi.Handle("GET", "/:id", &api.GetStreamServer{api.Common{DB:d, R:r}})

	streamParty := app.Party("/stream", sessionMiddleware)
	streamParty.Handle("GET", "/ws/publish/:id", &stream.PublishStream{api.Common{DB:d, R:r}})
	streamParty.Handle("GET", "/ws/subscribe/:id", &stream.WebSocketSubscriber{api.Common{DB:d, R:r}})
	streamParty.Handle("GET", "/http/subscribe/:id", &stream.SubscribeStream{api.Common{DB:d, R:r}})

	cricbuzzParty := app.Party("/api/cricbuzz")
	cricbuzzParty.Handle("GET", "/:id", &api.GetCricketScores{api.Common{DB:d, R:r}})

	app.Handle("GET", "/control", &control.Control{DB:d})

	if *debug{
		fmt.Println("enabling profling")
		app.Get("/debug/pprof/*action", pprof.New())
	}


	// if running in production mode listen on tls
	if *prod {
		fmt.Println("running in prod mode")
		cer, err := tls.LoadX509KeyPair(CERT_PATH, KEY_PATH)
 		if err != nil {
 			fmt.Println("unable to load keypair:", err)
 			return
 		}
 		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
		listener, err := tls.Listen("tcp4", "0.0.0.0:443", tlsConfig)
		app.Serve(listener)
	} else {
		fmt.Println("running in dev mode")
		app.Listen(":8000")
	}
}
