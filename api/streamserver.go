package api

import (
	"github.com/thekb/zyzz/db/models"
	"fmt"
	"gopkg.in/kataras/iris.v6"
)

type CreateStreamServer struct {
	Common
}

type GetStreamServer struct {
	Common
}

func (css *CreateStreamServer) Serve(ctx *iris.Context) {
	var streamServer models.StreamServer
	err := ctx.ReadJSON(&streamServer)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	streamServer.ShortId = getNewShortId()
	fmt.Println(streamServer)
	id, err := models.CreateStreamServer(css.DB, &streamServer)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	streamServer, _ = models.GetStreamServerForId(css.DB, id)
	ctx.JSON(iris.StatusOK, &streamServer)
	return
}

func (gss *GetStreamServer) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(SHORT_ID)
	streamServer, err := models.GetStreamServerForShortId(gss.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, &streamServer)
	return
}
