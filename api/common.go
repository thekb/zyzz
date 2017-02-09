package api

import (
	"net/http"
	"encoding/json"
)

type CommonApi struct {
}

func decodeRequestJSON(r *http.Request, destination interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(destination)
	return err
}
