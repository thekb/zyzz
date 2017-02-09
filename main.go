package main

import (
	"github.com/gorilla/mux"
	"github.com/thekb/zyzz/api"
	"net/http"
	"github.com/urfave/negroni"
	"fmt"
)

func main()  {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	streamRouter := apiRouter.PathPrefix("/stream").Subrouter()
	streamRouter.Handle("/", api.CreateStream{}).Methods("POST")
	streamRouter.Handle("/{id}", api.GetStream{}).Methods("GET")
	streamRouter.Handle("/{id}", api.PublishStream{}).Methods("PUT")
	n := negroni.Classic()
	n.UseHandler(r)
	fmt.Println("starting thanos...")
	http.ListenAndServe(":8000", n)
}
