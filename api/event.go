package api

import "gopkg.in/kataras/iris.v6"

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

func (ge GetEvent) Serve(ctx *iris.Context) {

}

func (ce CreateEvent) Serve(ctx *iris.Context) {

}

func (ge GetEvents) Serve(ctx *iris.Context) {

}

func (ges GetEventStreams) Serve(ctx *iris.Context) {

}

