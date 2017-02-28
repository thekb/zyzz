package control


import (
	"github.com/thekb/zyzz/api"
	"gopkg.in/kataras/iris.v6"
	ws "github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

type Control struct {
	api.Common
}

func (c *Control) Serve(ctx *iris.Context) {
	// verify if user is authenticated etc etc
	var err error
	var wsc *ws.Conn
	var msg []byte
	var upgrader = ws.Upgrader{
		CheckOrigin: func (r *http.Request) bool {return true},
	}

	wsc, err = upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, ctx.ResponseWriter.Header())
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &api.Response{Error:err.Error()})
		return
	}
	defer wsc.Close()
	for {
		_, msg, err = wsc.ReadMessage()
		if err != nil {
			fmt.Println("error received :", err)
			return
		}
		subscribe := HandleStreamMessage(msg)
		if subscribe {
			// run go routine to write messages
		}
	}



}

