package api

import (
	"net/http"
	"github.com/thekb/zyzz/db/models"
	"github.com/gorilla/mux"
)

type CreateUser struct {
	Common
}

type GetUser struct {
	Common
}

func (cuh CreateUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	cuh.decodeRequestJSON(r, &user)
	user.ShortId = getNewShortId()
	id, err := models.CreateUser(cuh.DB, &user)
	if err != nil {
		cuh.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
		return
	}
	user, _ = models.GetUserForId(cuh.DB, id)
	cuh.SendJSON(rw, &user)
}

func (guh GetUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortId := vars[SHORT_ID]
	user, err := models.GetUserForShortId(guh.DB, shortId)
	if err != nil {
		guh.SendErrorJSON(rw, err.Error(), http.StatusBadRequest)
	}
	guh.SendJSON(rw, &user)
}


