package api

import (
	"github.com/markbates/goth/providers/facebook"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/gothic"
	"github.com/markbates/goth"
	"github.com/thekb/zyzz/db/models"
	"fmt"
)

type FacebookCallback struct {
	Common
}


// ProviderIndex ...
type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func init() {
	goth.UseProviders(
		facebook.New("1359686897424077", "6acdd8af96d5d0ab7a83c023fda10049", "http://localhost:8000/auth/facebook/callback"),
	)
}

func Authenticate(ctx *iris.Context) {
	fmt.Println("authenticating")
	err := gothic.BeginAuthHandler(ctx)
	if err != nil {
		ctx.Log(iris.ProdMode, err.Error())
	}
}

func (fb *FacebookCallback) Serve(ctx *iris.Context) {
	user, err := gothic.CompleteUserAuth(ctx)
	fmt.Println("after authentication")
	if err != nil {
		ctx.SetStatusCode(iris.StatusUnauthorized)
		ctx.Writef(err.Error())
		return
	}
	user_model , err := models.GetUserForFBId(fb.DB, user.UserID)
	if err != nil {
		fmt.Println("User not present creating now:", err)
		var user_model models.User
		user_model.ShortId = getNewShortId()
		user_model.Email = user.Email
		user_model.Description = user.Description
		user_model.NickName = user.NickName
		user_model.FBId = user.UserID
		user_model.AvatarURL = user.AvatarURL
		user_model.AccessToken = user.AccessToken
		_, err := models.CreateUser(fb.DB, &user_model)
		if err != nil {
			ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
			return
		}
	}else {
		user_model.Email = user.Email
		user_model.Description = user.Description
		user_model.NickName = user.NickName
		user_model.AvatarURL = user.AvatarURL
		user_model.AccessToken = user.AccessToken
		models.UpdateUser(fb.DB, &user_model)
	}
	user_model, _ = models.GetUserForFBId(fb.DB, user.UserID)
	ctx.Session().Set("fbid", user_model.FBId)
	ctx.Session().Set("id", user_model.Id)
	ctx.Session().Set("short_id", user_model.ShortId)
	fmt.Println("Saved user is :", user_model)
	ctx.Redirect("/close")
	return
}
