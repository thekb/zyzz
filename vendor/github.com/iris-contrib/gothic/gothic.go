/*
This is a converted/porteed package in order to work with https://github.com/kataras/iris , from https://github.com/markbates/goth.

Package gothic wraps common behaviour when using Goth. This makes it quick, and easy, to get up
and running with Goth. Of course, if you want complete control over how things flow, in regards
to the authentication process, feel free and use Goth directly.

See https://github.com/iris-contrib/gothic/blob/master/_example/main.go & https://github.com/iris-contrib/gothic/blob/master/_example_low_level/main.go to see this in action.
*/
package gothic

import (
	"errors"
	"fmt"

	"gopkg.in/kataras/iris.v6"
	"github.com/markbates/goth"
)

// SessionName is the key used to access the session store.
// we could use the iris' sessions default, but this session should be not confict with the cookie session name defined by the sessions manager
const SessionName = "iris_gothic_session"

// IrisGothParams used to convert the context.URLParams to goth's params
type IrisGothParams map[string]string

// Get returns the value of
func (g IrisGothParams) Get(key string) string {
	return g[key]
}

var _ goth.Params = IrisGothParams{}

/*
BeginAuthHandler is a convienence handler for starting the authentication process.
It expects to be able to get the name of the provider from the named parameters
as either "provider" or url query parameter ":provider".

BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.

See https://github.com/iris-contrib/gothic/blob/master/_example_low_level/main.go to see this in action.
*/
func BeginAuthHandler(ctx *iris.Context) error {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.EmitError(400)
		return err
	}

	ctx.Redirect(url)
	return nil
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx *iris.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}

	return "state"

}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx *iris.Context) string {
	return ctx.URLParam("state")
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or url query parameter ":provider".

I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx *iris.Context) (string, error) {

	if ctx.Session() == nil {
		fmt.Println("You have to enable iris sessions")
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	ctx.Session().Set(SessionName, sess.Marshal())

	return url, nil
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.

It expects to be able to get the name of the provider from the named parameters
as either "provider" or url query parameter ":provider".

See https://github.com/iris-contrib/gothic/blob/master/_example_low_level/main.go to see this in action.
*/
var CompleteUserAuth = func(ctx *iris.Context) (goth.User, error) {

	if ctx.Session() == nil {
		fmt.Println("You have to enable iris sessions")
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	if ctx.Session().Get(SessionName) == nil {
		return goth.User{}, errors.New("could not find a matching session for this request")
	}

	sess, err := provider.UnmarshalSession(ctx.Session().GetString(SessionName))
	if err != nil {
		return goth.User{}, err
	}
	_, err = sess.Authorize(provider, IrisGothParams(ctx.URLParams()))

	if err != nil {
		return goth.User{}, err
	}

	return provider.FetchUser(sess)
}

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = getProviderName

func getProviderName(ctx *iris.Context) (string, error) {
	provider := ctx.Param("provider")
	if provider == "" {
		provider = ctx.URLParam(":provider")
	}
	if provider == "" {
		return provider, errors.New("you must select a provider")
	}
	return provider, nil
}
