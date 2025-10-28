package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/random"
	"github.com/songquanpeng/one-api/controller"
	"github.com/songquanpeng/one-api/model"
	"golang.org/x/oauth2"
)

func GetAuthOauth2Basic(ctx *gin.Context) {
	query := make(url.Values)
	query.Set("client_id", client_id)
	query.Set("response_type", "code")
	query.Set("redirect_uri", callback_url.String())
	state := random.GetUUID()
	query.Set("state", state)
	uri := url.URL{
		Scheme:   oauth2_url.Scheme,
		Host:     oauth2_url.Host,
		Path:     "/oauth/authorize",
		RawQuery: query.Encode(),
	}
	ctx.Redirect(302, uri.String())
}

func GetAuthOauth2BasicCallback(ctx *gin.Context) {
	if error_description := strings.TrimSpace(ctx.Query("error_description")); error_description != "" {
		ctx.JSON(http.StatusOK, ResponseJSON{
			Message: error_description,
		})
		return
	}
	code := strings.TrimSpace(ctx.Query("code"))
	token, err := getOauth2Token(code)
	if err != nil {
		ctx.JSON(http.StatusOK, ResponseJSON{
			Message: err.Error(),
		})
		return
	}
	userinfo, err := getOauth2Userinfo(token)
	if err != nil {
		ctx.JSON(http.StatusOK, ResponseJSON{
			Message: err.Error(),
		})
		return
	}
	user, err := insertKoolcenterUser(ctx, userinfo)
	if err != nil {
		ctx.JSON(http.StatusOK, ResponseJSON{
			Message: err.Error(),
		})
		return
	}
	if user.Status != model.UserStatusEnabled {
		ctx.JSON(http.StatusOK, ResponseJSON{
			Message: "用户已被封禁",
		})
		return
	}
	// ctx.JSON(http.StatusOK, user)
	// controller.SetupLogin(user, ctx)
	controller.SetCookie(user, ctx)
	ctx.Data(200, "text/html; charset=utf-8", []byte(showHTML(user)))

}

func getOauth2Token(code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("not foun code")
	}
	values := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	values.Set("client_id", client_id)
	values.Set("client_secret", client_secret)
	values.Set("redirect_uri", callback_url.String())
	req, err := http.NewRequest(http.MethodPost, oauth2_url.JoinPath("/oauth/token").String(), strings.NewReader(values.Encode()))
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
	uri := fmt.Sprintf("%v/oauth/userinfo?access_token=%s", oauth2_url.String(), token.AccessToken)
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
		user.Username = "kc_" + strconv.Itoa(model.GetMaxUserId()+1)
		user.Email = userinfo.Email
		user.DisplayName = userinfo.Name
		if user.DisplayName == "" {
			user.DisplayName = "KoolcenterUser"
		}
		user.Role = model.RoleCommonUser
		user.Status = model.UserStatusEnabled
		// 默认额度
		// user.Quota = 10000
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

type ResponseJSON struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func showHTML(user *model.User) string {
	by, _ := json.Marshal(user)
	obj := string(by)
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title></title>
</head>
<body>
        <script>
            localStorage.setItem('user', '` + obj + `');
			location.href = "/"
        </script>
</body>
</html>
`
	return html
}
