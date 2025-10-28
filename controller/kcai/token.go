package kcai

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/config"
	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/controller/oauth2"
	"github.com/songquanpeng/one-api/model"
	"github.com/songquanpeng/one-api/relay/adaptor/lightrag"
)

type ResponseKcai struct {
	Token string `json:"token"`
	Model string `json:"model"`
}

func getUserToken(userid int) (string, error) {
	tokens, err := model.GetAllUserTokens(userid, 0, config.ItemsPerPage, "")
	if err != nil || tokens == nil || len(tokens) == 0 {
		return "", errors.New("not found user token")
	}
	var token string
	for _, v := range tokens {
		if v.Status == 1 {
			token = v.Key
			break
		}
	}
	if token == "" {
		return "", errors.New("not found token")
	}
	return "sk-" + token, nil
}
func GetKcaiToken(ctx *gin.Context) {
	userId := ctx.GetInt(ctxkey.Id)
	token, err := getUserToken(userId)
	if err != nil {
		ctx.JSON(http.StatusOK, oauth2.ResponseJSON{
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, oauth2.ResponseJSON{
		Data: ResponseKcai{
			Token: token,
			Model: strings.Join(lightrag.ModelList, ","),
		},
	})
}
