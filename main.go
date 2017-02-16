package main

import (
	"github.com/gorilla/mux"
	"github.com/thekb/zyzz/api"
	"net/http"
	"github.com/urfave/negroni"
	"fmt"
	"github.com/thekb/zyzz/db"
	"github.com/thekb/zyzz/stream"
	"flag"
)

const (
	CERT_PATH = "/etc/letsencrypt/live/zyzz.co/fullchain.pem"
	KEY_PATH = "/etc/letsencrypt/live/zyzz.co/privkey.pem"
)

var tlsFlag = flag.Bool("tls", false, "enable tls")

func main() {
	flag.Parse()

	d, err := db.GetDB()
	if err != nil {
		fmt.Println("unable to connect to db:", err)
		return
	}

	r := mux.NewRouter()
	// register API methods
	apiRouter := r.PathPrefix("/api").Subrouter().StrictSlash(true)
	apiRouter.Handle("/user/", api.CreateUser{api.Common{DB:d}}).Methods("POST")
	apiRouter.Handle("/user/{shortId}/", api.GetUser{api.Common{DB:d}}).Methods("GET")
	apiRouter.Handle("/streamserver/", api.CreateStreamServer{api.Common{DB:d}}).Methods("POST")
	apiRouter.Handle("/streamserver/{shortId}/", api.GetStreamServer{api.Common{DB:d}}).Methods("GET")
	apiRouter.Handle("/stream/", api.GetStreams{api.Common{DB:d}}).Methods("GET")
	apiRouter.Handle("/stream/", api.CreateStream{api.Common{DB:d}}).Methods("POST")
	apiRouter.Handle("/stream/{shortId}/", api.GetStream{api.Common{DB:d}}).Methods("GET")

	// register stream methods
	streamRouter := r.PathPrefix("/stream").Subrouter()
	streamRouter.Handle("/publish/{shortId}/", stream.PublishStream{api.Common{DB:d}}).Methods("POST")
	streamRouter.Handle("/subscribe/{shortId}/", stream.SubscribeStream{api.Common{DB:d}}).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(r)
	fmt.Println("starting zyzz...")
	if *tlsFlag {
		fmt.Println("using tls...")
		http.ListenAndServeTLS(":443", CERT_PATH, KEY_PATH, n)
	} else {
		http.ListenAndServe(":8000", n)

	}
}
