package api

import (
	"net/http"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/ventu-io/go-shortid"
)

const (
	HEADER_CONTENT_TYPE = "Content-Type"
	JSON_CONTENT_TYPE = "application/json"
	SHORT_ID = "id"
)

var shortIdGenerator *shortid.Shortid

func init () {
	initializeShortId()
}

type Common struct {
	DB *sqlx.DB
}

type Response struct {
	Error string `json:"error"`
	Data interface{} `json:"data"`
}

func (c *Common) decodeRequestJSON(r *http.Request, destination interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(destination)
	return err
}

func (c *Common) SendErrorJSON(rw http.ResponseWriter, error string, status int) {
	rw.Header().Set(HEADER_CONTENT_TYPE, JSON_CONTENT_TYPE)
	rw.WriteHeader(status)
	response := Response{Error: error}
	b, _ := json.Marshal(&response)
	rw.Write(b)
	return
}

func (c *Common) SendJSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set(HEADER_CONTENT_TYPE, JSON_CONTENT_TYPE)
	rw.WriteHeader(http.StatusOK)
	response := Response{Data: data}
	b, _ := json.Marshal(&response)
	rw.Write(b)
}

func initializeShortId() {
	shortIdGenerator, _ = shortid.New(1, shortid.DefaultABC, 1729)
}

func getNewShortId() string {
	shortId, _ := shortIdGenerator.Generate()
	return shortId
}