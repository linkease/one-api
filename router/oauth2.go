package router

import (
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/controller/kcai"
	"github.com/songquanpeng/one-api/controller/oauth2"
	"github.com/songquanpeng/one-api/middleware"
)

// http://localhost:3000/auth/oauth2_basic
// http://localhost:9096/oauth/authorize
func HandleOauth2Router(router *gin.Engine) {
	router.GET("/api/kcai/token", middleware.UserAuth(), kcai.GetKcaiToken)
	router.GET("/api/kcai/nextchat", kcai.GetKcaiNextchat)
	router.GET("/auth/oauth2_basic", oauth2.GetAuthOauth2Basic)
	router.GET("/auth/oauth2_basic/callback", oauth2.GetAuthOauth2BasicCallback)
}
