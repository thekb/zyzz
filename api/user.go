package api

import (
	"github.com/thekb/zyzz/db/models"
	"gopkg.in/kataras/iris.v6"
	"fmt"
)

type CreateUser struct {
	Common
}

type GetUser struct {
	Common
}

type GetUserStream struct {
	Common
}

type GetCurrentUserStream struct {
	Common
}

type UpdateUser struct {
	Common
}

func (cuh CreateUser) Serve(ctx *iris.Context)  {
	var user models.User
	err := ctx.ReadJSON(&user)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to decode request %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	user.ShortId = getNewShortId()
	id, err := models.CreateUser(cuh.DB, &user)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	user, _ = models.GetUserForId(cuh.DB, id)
	ctx.JSON(iris.StatusOK, Response{Data:user})
	return
}

func (guh GetUser) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(SHORT_ID)
	user, err := models.GetUserForShortId(guh.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, Response{Data:user})
	return
}

func (uu UpdateUser) Serve(ctx *iris.Context) {
	shortId := ctx.GetString(SHORT_ID)
	user, err := models.GetUserForShortId(uu.DB, shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	err = ctx.ReadJSON(&user)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to decode request %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	err = models.UpdateUser(uu.DB, &user)
	if err != nil {
		ctx.Log(iris.ProdMode, "unable to update user %s", err.Error())
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	ctx.JSON(iris.StatusOK, Response{Data:user})
	return
}

func (gus GetUserStream) Serve(ctx *iris.Context) {

}

func (gcus GetCurrentUserStream) Serve(ctx *iris.Context) {
	user_shortId := ctx.GetString(SHORT_ID)
	stream, err := models.GetUserActiveStream(gcus.DB, user_shortId)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, &Response{Error:err.Error()})
		return
	}
	user, err  := models.GetUserForId(gcus.DB, int64(stream.CreatorId))
	if err != nil {
		fmt.Println("unable to get user for stream ", err)
	}
	stream.User = user
	ctx.JSON(iris.StatusOK, &Response{Data:stream})
	return
}