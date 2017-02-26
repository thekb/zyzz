package api

import (
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/db/models"
)

type GetEvent struct {
	Common
}

type GetEvents struct {
	Common
}

type CreateEvent struct {
	Common
}

type GetEventStreams struct {
	Common
}

func (ge *GetEvent) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(SHORT_ID)
	event, err := models.GetEventForShortId(ge.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, Response{Data:event})
	return
}

func (ce *CreateEvent) Serve(ctx *iris.Context) {
	var event models.Event
	err := ctx.ReadJSON(&event)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to decode request %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	event.ShortId = getNewShortId()
	id, err := models.CreateEvent(ce.DB, &event)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	event, _ = models.GetEventForId(ce.DB, id)
	ctx.JSON(iris.StatusOK, Response{Data:event})
	return
}

func (ge *GetEvents) Serve(ctx *iris.Context) {

}

func (ges *GetEventStreams) Serve(ctx *iris.Context) {

}

