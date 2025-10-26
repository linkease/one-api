package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/random"
	"github.com/songquanpeng/one-api/controller"
	"github.com/songquanpeng/one-api/middleware"
	"github.com/songquanpeng/one-api/model"
	"golang.org/x/oauth2"
)

const (
	client_id     = "koolcenter"
	client_secret = "123456"
	callback_url  = "http://localhost:3000/auth/oauth2_basic/callback"
)

var (
	target_url = url.URL{
		Scheme: "http",
		Host:   "localhost:9096",
	}
)

// http://localhost:3000/auth/oauth2_basic
// http://localhost:9096/oauth/authorize
func HandleOauth2Router(router *gin.Engine) {
	router.GET("/auth/oauth2_basic", middleware.CriticalRateLimit(), GetAuthOauth2Basic)
	router.GET("/auth/oauth2_basic/callback", middleware.CriticalRateLimit(), GetAuthOauth2BasicCallback)
}

func GetAuthOauth2Basic(ctx *gin.Context) {
	query := make(url.Values)
	query.Set("client_id", client_id)
	query.Set("response_type", "code")
	query.Set("redirect_uri", callback_url)
	state := random.GetUUID()
	query.Set("state", state)
	uri := url.URL{
		// Scheme:   "https",
		// Host:     "sso.koolcenter.com",
		Scheme:   target_url.Scheme,
		Host:     target_url.Host,
		Path:     "/oauth/authorize",
		RawQuery: query.Encode(),
	}
	ctx.Redirect(302, uri.String())
}

func GetAuthOauth2BasicCallback(ctx *gin.Context) {
	code := strings.TrimSpace(ctx.Query("code"))
	token, err := getOauth2Token(code)
	if err != nil {
		ctx.JSON(http.StatusOK, ResponseError{
			Message: err.Error(),
		})
		return
	}
	userinfo, err := getOauth2Userinfo(token)
	if err != nil {
		ctx.JSON(http.StatusOK, ResponseError{
			Message: err.Error(),
		})
		return
	}
	user, err := insertKoolcenterUser(ctx, userinfo)
	if err != nil {
		ctx.JSON(http.StatusOK, ResponseError{
			Message: err.Error(),
		})
		return
	}
	if user.Status != model.UserStatusEnabled {
		ctx.JSON(http.StatusOK, ResponseError{
			Message: "用户已被封禁",
		})
		return
	}
	// ctx.JSON(http.StatusOK, user)
	controller.SetupLogin(user, ctx)
}

func getOauth2Token(code string) (*oauth2.Token, error) {
	values := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	values.Set("client_id", client_id)
	values.Set("client_secret", client_secret)
	values.Set("redirect_uri", callback_url)
	req, err := http.NewRequest("POST", target_url.JoinPath("/oauth/token").String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(url.QueryEscape(client_id), url.QueryEscape(client_secret))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var token oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}
func getOauth2Userinfo(token *oauth2.Token) (*Userinfo, error) {
	uri := fmt.Sprintf("%v/oauth/userinfo?access_token=%s", target_url.String(), token.AccessToken)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userinfo Userinfo
	if err := json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		return nil, err
	}
	return &userinfo, nil
}
func insertKoolcenterUser(ctx context.Context, userinfo *Userinfo) (*model.User, error) {
	user := model.User{
		KcId: userinfo.Uid,
	}
	var err error
	// 已经注册了
	if model.IsKoolcenterIdAlreadyTaken(userinfo.Uid) {
		err = user.FillUserByKcId()
	} else {
		user.Username = "kc_" + random.GetUUID()
		user.Email = userinfo.Email
		user.DisplayName = userinfo.Name
		if user.DisplayName == "" {
			user.DisplayName = "KoolcenterUser"
		}
		user.Role = model.RoleCommonUser
		user.Status = model.UserStatusEnabled
		err = user.Insert(ctx, 0)
	}
	if err != nil {
		return nil, err
	}
	if user.Status != model.UserStatusEnabled {
		return nil, errors.New("用户已被封禁")
	}
	return &user, nil
}

type Userinfo struct {
	Uid      string `json:"uid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResponseError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
