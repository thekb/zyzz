package api

import (
	"github.com/thekb/zyzz/db/models"
	"fmt"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/control"
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

	//err := ctx.ReadJSON(&stream)
	//if err != nil {
	//	ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
	//	return
	//}
	event_shortId := ctx.GetString(SHORT_ID)
	stream.ShortId = getNewShortId()
	stream.EventId = event_shortId
	defaultStreamServer := models.GetDefaultStreamServer(cs.DB)
	user_id, _ := ctx.Session().GetInt("id")
	stream.CreatorId = user_id
	stream.StreamServerId = defaultStreamServer.Id
	stream.TransportUrl = fmt.Sprintf(TRANSPORT_URL_FORMAT, stream.ShortId)
	stream.PublishUrl = fmt.Sprintf(PUBLISH_URL_FORMAT, defaultStreamServer.HostName, stream.ShortId)
	stream.SubscribeUrl = fmt.Sprintf(SUBSCRIBE_URL_FORMAT, defaultStreamServer.HostName, stream.ShortId)
	id, err := models.CreateStream(cs.DB, &stream)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	// setup stream sockets
	err = control.StreamMap.CreateStream(stream.ShortId, stream.CreatorId)
	if err != nil {
		fmt.Println("unable to setup stream:", err)
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
	ctx.JSON(iris.StatusOK, &Response{Data:stream})
	return
}

func (gs *GetStreams) Serve(ctx *iris.Context) {
	event_shortId := ctx.GetString(SHORT_ID)
	streams, err := models.GetStreams(gs.DB, event_shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, &Response{Data:streams})
	return
}
