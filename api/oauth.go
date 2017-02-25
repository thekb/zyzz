package api

import (
	"github.com/markbates/goth/providers/facebook"
	"github.com/iris-contrib/gothic"
	"github.com/markbates/goth"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/zyzz/db/models"
)

type Callback struct {
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
	err := gothic.BeginAuthHandler(ctx)
	if err != nil {
		ctx.Log(iris.ProdMode, err.Error())
	}
}

func (cb Callback) Serve(ctx *iris.Context) {
	user, err := gothic.CompleteUserAuth(ctx)
	if err != nil {
		ctx.SetStatusCode(iris.StatusUnauthorized)
		ctx.Writef(err.Error())
		return
	}
	var user_model models.User
	user_model.ShortId = getNewShortId()
	user_model.Email = user.Email
	user_model.Description = user.Description
	user_model.NickName = user.NickName
	user_model.FBId = user.UserID
	user_model.AvatarURL = user.AvatarURL
	id, err := models.CreateUser(cb.DB, &user_model)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, Response{Error:err.Error()})
		return
	}
	user_model, _ = models.GetUserForId(cb.DB, id)
	ctx.JSON(iris.StatusOK, Response{Data:user})
	return

}
