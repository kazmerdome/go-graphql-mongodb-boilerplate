package service

import (
	"aery-graphql/config"
	"aery-graphql/guard"
	"aery-graphql/model"
	"aery-graphql/utility"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthController ...
type AuthController struct{}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	FirstName     string `json:"given_name"`
	Lastname      string `json:"family_name"`
}

type responseToken struct {
	Token string `token:"id"`
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.GetSecret("GOOGLE_CALLBACK_URL"),
		ClientID:     config.GetSecret("GOOGLE_CLIENT_ID"),
		ClientSecret: config.GetSecret("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

func generateStateOauthCookie(c echo.Context) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 32)
	rand.Read(b)
	// state := base64.URLEncoding.EncodeToString(b)
	state := "BArEYzvg_UDx_Y1kx-DWCMEZP3WWJUze5tkF4NH1ZyU="
	cookie := http.Cookie{Name: "aerylabs-tmpst", Value: state, Expires: expiration, Path: "/"}
	c.SetCookie(&cookie)
	return state
}

// OauthGoogle ...
func (controller AuthController) OauthGoogle(c echo.Context) error {
	oauthState := generateStateOauthCookie(c)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

// OauthGoogleCallback ...
func (controller AuthController) OauthGoogleCallback(c echo.Context) error {
	oauthState, err := c.Cookie("aerylabs-tmpst")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	qState := c.QueryParam("state")
	if qState != oauthState.Value {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid oauth google state")
	}

	qCode := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, qCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response, err := http.Get(googleUserInfoURL + token.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)

	userInfo := googleUserInfo{}
	if err := json.Unmarshal(content, &userInfo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !userInfo.VerifiedEmail {
		return echo.NewHTTPError(http.StatusInternalServerError, "unverified google email")
	}

	user := model.User{}
	filter := model.UserWhereInput{Email: &userInfo.Email}
	if err := user.One(&filter); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ADD GOOGLE TO EXISTING ACCOUNT | if user exist but does not have googleID
	if user.Email == userInfo.Email && user.GoogleID != userInfo.ID {
		if user.GoogleID != "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "mismatch error, user already has googleID")
		}
		user.GoogleID = userInfo.ID
		if err := user.Update(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// SIGNUP | if user does not exist
	if user.Email == "" {
		user.Firstname = userInfo.FirstName
		user.Lastname = userInfo.Lastname
		user.Email = userInfo.Email
		user.AppPolicy = []model.AppPolicy{model.AppPolicy{Resource: model.BlackdomeServer, Role: guard.User}}
		user.GoogleID = userInfo.ID
		if err := user.Create(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// LOGIN | if user exist and googeID exist
	if user.Email == userInfo.Email && user.GoogleID == userInfo.ID {
		jwt, err := utility.GenerateJWTToken(user.ID.Hex())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		response := responseToken{Token: jwt}
		return c.JSON(http.StatusOK, response)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "authentication error")
}
