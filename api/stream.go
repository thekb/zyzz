package api

import (
	"net/http"
	"github.com/thekb/zyzz/db/models"
	"fmt"
	"github.com/gorilla/mux"
)

type CreateStream struct {
	Common
}

type GetStream struct {
	Common
}

type GetStreams struct {
	Common
}

const (
	TRANSPORT_URL_FORMAT = "ipc:///tmp/stream_%s.ipc"
	ENDPOINT_URL_FORMAT = "https://%s/stream/%s/"

)

func (cs CreateStream) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var stream models.Stream
	cs.decodeRequestJSON(r, &stream)
	stream.ShortId = getNewShortId()
	defaultStreamServer := models.GetDefaultStreamServer(cs.DB)
	stream.CreatorId = models.GetDefaultUser(cs.DB).Id
	stream.StreamServerId = defaultStreamServer.Id
	stream.TransportUrl = fmt.Sprintf(TRANSPORT_URL_FORMAT, stream.ShortId)
	stream.EndPoint = fmt.Sprintf(ENDPOINT_URL_FORMAT, defaultStreamServer.HostName, stream.ShortId)
	id, err := models.CreateStream(cs.DB, &stream)
	if err != nil {
		cs.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
	}
	stream, _ = models.GetStreamForId(cs.DB, id)
	cs.SendJSON(rw, &stream)
	return
}

func (gs GetStream) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortId := vars[SHORT_ID]

	stream, err := models.GetStreamForShortId(gs.DB, shortId)
	if err != nil {
		gs.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
	}
	gs.SendJSON(rw, &stream)
	return
}

func (gs GetStreams) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	streams, err := models.GetStreams(gs.DB)
	if err != nil {
		gs.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
		return
	}
	gs.SendJSON(rw, &streams)
	return
}
