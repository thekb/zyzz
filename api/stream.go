package api

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thekb/zyzz/stream"
	"encoding/json"
)

type GetStream struct {
	*CommonApi
}
type CreateStream struct {
	*CommonApi
}
type PublishStream struct {
	*CommonApi
}

type CreateStreamRequest struct {
	Name string `json:"name"`
}

type CreateStreamResponse struct {
	Id string `json:"id"`
}

func (gs GetStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "audio/aac")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	err := stream.SubscribeStream(id, w)
	fmt.Printf("get stream error %s", err)

}

func (cs CreateStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var createStreamRequest CreateStreamRequest
	err := decodeRequestJSON(r, &createStreamRequest)
	if err != nil {

	}
	fmt.Printf("creating stream with name %s", createStreamRequest.Name)
	id := stream.CreateStream(createStreamRequest.Name)
	createStreamResponse := CreateStreamResponse{}
	createStreamResponse.Id = id
	data, _ := json.Marshal(createStreamResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)



}

func (ps PublishStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := stream.PublishStream(id, r.Body)
	if err != nil {
		fmt.Printf("publish error %s", err)
	}

}