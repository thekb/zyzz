package api

import (
	"github.com/thekb/zyzz/db/models"
	"fmt"
	"gopkg.in/kataras/iris.v6"
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
	PUBLISH_URL_FORMAT = "https://%s/stream/publish/%s/"
	SUBSCRIBE_URL_FORMAT = "https://%s/stream/subscribe/%s/"

)

func (cs *CreateStream) Serve(ctx *iris.Context) {
	var stream models.Stream
	err := ctx.ReadJSON(&stream)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	stream.ShortId = getNewShortId()
	defaultStreamServer := models.GetDefaultStreamServer(cs.DB)
	stream.CreatorId = models.GetDefaultUser(cs.DB).Id
	stream.StreamServerId = defaultStreamServer.Id
	stream.TransportUrl = fmt.Sprintf(TRANSPORT_URL_FORMAT, stream.ShortId)
	stream.PublishUrl = fmt.Sprintf(PUBLISH_URL_FORMAT, defaultStreamServer.HostName, stream.ShortId)
	stream.SubscribeUrl = fmt.Sprintf(SUBSCRIBE_URL_FORMAT, defaultStreamServer.HostName, stream.ShortId)
	id, err := models.CreateStream(cs.DB, &stream)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	stream, _ = models.GetStreamForId(cs.DB, id)
	ctx.JSON(iris.StatusOK, &stream)
	return
}

func (gs *GetStream) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(SHORT_ID)

	stream, err := models.GetStreamForShortId(gs.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, &stream)
	return
}

func (gs *GetStreams) Serve(ctx *iris.Context) {
	streams, err := models.GetStreams(gs.DB)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, &streams)
	return
}
