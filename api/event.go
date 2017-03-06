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

type UpdateEvent struct {
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

func (ce *UpdateEvent) Serve(ctx *iris.Context) {
	//var event models.Event
	event_shortId := ctx.GetString(SHORT_ID)
	event, err := models.GetEventForShortId(ce.DB, event_shortId)
	if err != nil {
		ctx.Log(iris.ProdMode, "Not able to fetch an event %s", err.Error())
		ctx.JSON(iris.StatusNotFound, Response{Error:err.Error()})
		return
	}
	err = ctx.ReadJSON(&event)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to decode request %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	err = models.UpdateEvent(ce.DB, &event)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, Response{Error:err.Error()})
		return
	}
	event, _ = models.GetEventForId(ce.DB, event.Id)
	ctx.JSON(iris.StatusOK, Response{Data:event})
	return
}

func (ge *GetEvents) Serve(ctx *iris.Context) {
	events, err := models.GetEvents(ge.DB)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, &Response{Data:events})
	return
}

func (ges *GetEventStreams) Serve(ctx *iris.Context) {
	event_shortId := ctx.GetString(SHORT_ID)
	streams, err := models.GetStreams(ges.DB, event_shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	for i := 0; i < len(streams); i++ {
		user_publish, err := models.GetUserForId(ges.DB, int64(streams[i].CreatorId))
		if err == nil {
			streams[i].User = user_publish
		}

	}
	ctx.JSON(iris.StatusOK, &Response{Data:streams})
	return
}

