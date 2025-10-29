package kcai

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/configs"
	"github.com/songquanpeng/one-api/controller/oauth2"
)

func GetKcaiNextchat(ctx *gin.Context) {
	userId := ctx.GetInt(ctxkey.Id)
	// 未登陆
	if userId <= 0 {
		ctx.Redirect(302, "/auth/oauth2_basic")
		return
	}
	token, err := getUserToken(userId)
	if err != nil {
		ctx.JSON(http.StatusOK, oauth2.ResponseJSON{
			Message: err.Error(),
		})
		return
	}
	nextchatUrl := configs.KcaiConfig.GetNextchatUrl() + `?token=` + token
	ctx.Redirect(302, nextchatUrl)
}
