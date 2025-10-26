package router

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/config"
	"github.com/songquanpeng/one-api/controller"
	"github.com/songquanpeng/one-api/middleware"
)

func SetWebRouter(router *gin.Engine, buildFS embed.FS) {
	indexPageData, _ := buildFS.ReadFile(fmt.Sprintf("web/build/%s/index.html", config.Theme))

	// Create a router group for the admin interface
	webRouter := router.Group("/kc-admin")
	webRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	webRouter.Use(middleware.GlobalWebRateLimit())
	webRouter.Use(middleware.Cache())
	webRouter.Use(static.Serve("/kc-admin", common.EmbedFolder(buildFS, fmt.Sprintf("web/build/%s", config.Theme))))

	// Redirect root to /kc-admin/
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/kc-admin/")
	})

	router.NoRoute(func(c *gin.Context) {
		// Only /v1 routes are handled specially for OpenAI API compatibility
		if strings.HasPrefix(c.Request.RequestURI, "/v1") {
			controller.RelayNotFound(c)
			return
		}
		// All other routes serve the SPA (including /kc-admin/*)
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPageData)
	})
}
