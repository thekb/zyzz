package api

import (
	"github.com/markbates/goth/providers/facebook"
	"gopkg.in/kataras/iris.v6"
	"github.com/thekb/gothic"
	"github.com/markbates/goth"
	"github.com/thekb/zyzz/db/models"
	"fmt"
	"github.com/markbates/goth/providers/gplus"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"os"
	"net/http"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"io"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

const (
	FACEBOOK_END_POINT = "https://graph.facebook.com/me?fields=email,first_name,last_name,link,about,id,name,picture,location"
	FACEBOOK_CALLBACK_URL_FORMAT = "%s/auth/facebook/callback"
	GPLUS_CALLBACK_URL_FORMAT = "%s/auth/gplus/callback"
	GPLUS_END_POINT = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type AuthCallback struct {
	Common
}

type AppTokenVerify struct {
	Common
}

type TokenInfo struct {
	AccessToken string `json:"access_token"`
	Provider    string `json:"provider"`
	User        models.User `json:"user"`
}

// ProviderIndex ...
type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func init() {
	goth.UseProviders(
		facebook.New(os.Getenv("FacebookKey"), os.Getenv("FacebookSecret"), fmt.Sprintf(FACEBOOK_CALLBACK_URL_FORMAT, os.Getenv("AppServer"))),
		gplus.New(os.Getenv("GplusKey"), os.Getenv("GplusSecret"), fmt.Sprintf(GPLUS_CALLBACK_URL_FORMAT, os.Getenv("AppServer"))),
		//facebook.New("1359686897424077", "6acdd8af96d5d0ab7a83c023fda10049","https://shortpitch.live/auth/facebook/callback"),
		//gplus.New("882975999128-sjon5nvc0vsoi71ia3u15qg0hqkmc209.apps.googleusercontent.com","zubWjfuROIyXN8KNtBKFW62E", "https://www.shortpitch.live/auth/gplus/callback/"),
	)
}

func Authenticate(ctx *iris.Context) {
	err := gothic.BeginAuthHandler(ctx)
	if err != nil {
		ctx.Log(iris.ProdMode, err.Error())
	}
}

func (ac *AuthCallback) Serve(ctx *iris.Context) {
	gothUser, err := gothic.CompleteUserAuth(ctx)
	if err != nil {
		ctx.SetStatusCode(iris.StatusUnauthorized)
		ctx.Writef(err.Error())
		return
	}
	user, err := createOrUpdateUser(ac.DB, gothUser)
	ctx.Session().Set("fbid", user.FBId)
	ctx.Session().Set("id", user.Id)
	ctx.Session().Set("short_id", user.ShortId)
	fmt.Println("Saved user is :", user)
	ctx.Redirect("/close")
	return
}

func (atv *AppTokenVerify) Serve(ctx *iris.Context) {
	var tokenInfo TokenInfo
	err := ctx.ReadJSON(&tokenInfo)
	if err != nil {
		ctx.SetStatusCode(iris.StatusBadRequest)
		ctx.Writef(err.Error())
		return
	}
	switch tokenInfo.Provider {
	case "facebook":
		keyHash := hmac.New(sha256.New, []byte(os.Getenv("FacebookSecret")))
		keyHash.Write([]byte(tokenInfo.AccessToken))
		appsecretProof := hex.EncodeToString(keyHash.Sum(nil))
		req, _ := http.NewRequest("GET", FACEBOOK_END_POINT + "&access_token=" + url.QueryEscape(tokenInfo.AccessToken) + "&appsecret_proof=" + appsecretProof, bytes.NewBufferString(""))
		client := &http.Client{}
		resp, err := client.Do(req)
		// close response body
		defer resp.Body.Close()
		if err != nil {
			ctx.SetStatusCode(iris.StatusUnauthorized)
			ctx.Writef(err.Error())
			return
		}
		if resp.StatusCode == http.StatusOK {
			var gothUser goth.User
			bits, err := ioutil.ReadAll(resp.Body)
			err = json.NewDecoder(bytes.NewReader(bits)).Decode(&gothUser.RawData)
			if err != nil {
				ctx.SetStatusCode(iris.StatusUnauthorized)
				ctx.Writef(err.Error())
				return
			}
			err = fbUserFromReader(bytes.NewReader(bits), &gothUser)
			if err != nil {
				ctx.SetStatusCode(iris.StatusUnauthorized)
				ctx.Writef(err.Error())
				return
			}
			gothUser.AccessToken = tokenInfo.AccessToken
			user, err := createOrUpdateUser(atv.DB, gothUser)
			if err != nil {
				fmt.Println("error creating or updating user", err)
			}
			uuid4 := uuid.NewV4()
			tokenInfo.AccessToken = uuid4.String()
			tokenInfo.User = user
			atv.R.Set(uuid4.String(), 1, 0)
			ctx.JSON(iris.StatusOK, Response{Data:tokenInfo})
			return
			//atv.R.Set()
		} else {
			ctx.SetStatusCode(iris.StatusUnauthorized)
			ctx.Writef(err.Error())
			return
		}

	case "gplus":
		req, _ := http.NewRequest("GET", GPLUS_END_POINT + "?access_token=" + url.QueryEscape(tokenInfo.AccessToken), bytes.NewBufferString(""))
		fmt.Println(GPLUS_END_POINT + "?access_token=" + url.QueryEscape(tokenInfo.AccessToken))
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			ctx.SetStatusCode(iris.StatusUnauthorized)
			ctx.Writef(err.Error())
			return
		}
		fmt.Println(resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			var gothUser goth.User
			bits, err := ioutil.ReadAll(resp.Body)
			err = json.NewDecoder(bytes.NewReader(bits)).Decode(&gothUser.RawData)
			if err != nil {
				ctx.SetStatusCode(iris.StatusUnauthorized)
				ctx.Writef(err.Error())
				return
			}
			err = gplusUserFromReader(bytes.NewReader(bits), &gothUser)
			if err != nil {
				ctx.SetStatusCode(iris.StatusUnauthorized)
				ctx.Writef(err.Error())
				return
			}
			gothUser.AccessToken = tokenInfo.AccessToken
			user, err := createOrUpdateUser(atv.DB, gothUser)
			if err != nil {
				fmt.Println("error creating or updating user", err)
			}
			uuid4 := uuid.NewV4()
			tokenInfo.AccessToken = uuid4.String()
			tokenInfo.User = user
			atv.R.Set(uuid4.String(), 1, 0)
			ctx.JSON(iris.StatusOK, Response{Data:tokenInfo})
			return
		}
	default:
		ctx.SetStatusCode(iris.StatusBadRequest)
		ctx.Writef(err.Error())
		return
	}
}

func createOrUpdateUser(DB *sqlx.DB, user goth.User) (models.User, error) {
	user_model, err := models.GetUserForFBId(DB, user.UserID)
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
		_, err := models.CreateUser(DB, &user_model)
		if err != nil {
			return user_model, err
		}
	} else {
		fmt.Println("updating user")
		user_model.Email = user.Email
		user_model.Description = user.Description
		user_model.NickName = user.NickName
		user_model.AvatarURL = user.AvatarURL
		user_model.AccessToken = user.AccessToken
		models.UpdateUser(DB, &user_model)
	}
	user_model, err = models.GetUserForFBId(DB, user.UserID)
	return user_model, err
}

func fbUserFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		About     string `json:"about"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Link      string `json:"link"`
		Picture   struct {
				  Data struct {
					       URL string `json:"url"`
				       } `json:"data"`
			  } `json:"picture"`
		Location  struct {
				  Name string `json:"name"`
			  } `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.Name
	user.Email = u.Email
	user.Description = u.About
	user.AvatarURL = u.Picture.Data.URL
	user.UserID = u.ID
	user.Location = u.Location.Name

	return err
}


func gplusUserFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Link      string `json:"link"`
		Picture   string `json:"picture"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.Name
	user.Email = u.Email
	//user.Description = u.Bio
	user.AvatarURL = u.Picture
	user.UserID = u.ID
	//user.Location = u.Location.Name

	return err
}