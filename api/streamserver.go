package api

import (
	"net/http"
	"github.com/thekb/zyzz/db/models"
	"github.com/gorilla/mux"
	"fmt"
)

type CreateStreamServer struct {
	Common
}

type GetStreamServer struct {
	Common
}

func (css CreateStreamServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var streamServer models.StreamServer
	css.decodeRequestJSON(r, &streamServer)
	streamServer.ShortId = getNewShortId()
	fmt.Println(streamServer)
	id, err := models.CreateStreamServer(css.DB, &streamServer)
	if err != nil {
		css.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
	}
	streamServer, _ = models.GetStreamServerForId(css.DB, id)
	css.SendJSON(rw, &streamServer)
	return
}

func (gss GetStreamServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortId := vars[SHORT_ID]

	streamServer, err := models.GetStreamServerForShortId(gss.DB, shortId)
	if err != nil {
		gss.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
		return
	}
	gss.SendJSON(rw, &streamServer)
	return
}
