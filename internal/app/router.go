package app

import (
	"embed"
	"net/http"
	"time"

	"github.com/Lofanmi/gobana/service"
	"github.com/gin-gonic/gin"
)

//go:embed dist/index.html
var distIndexHTML string

//go:embed dist/favicon.ico
var distFaviconIcon []byte

//go:embed dist/static/*
var distStaticFS embed.FS

func (s *Application) registerStaticRouter(engine *gin.Engine) {
	engine.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", distFaviconIcon)
	})
	if s.ConfigApplication.Production {
		engine.GET("/", func(c *gin.Context) {
			c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
			c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
			c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, distIndexHTML)
		})
		engine.StaticFS("/frontend", http.FS(distStaticFS))
	}
}

func (s *Application) registerApiRouter(router gin.IRoutes) {
	router.GET("config/backend_list", func(c *gin.Context) {
		var req service.GetBackendListRequest
		if err := c.Bind(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		resp, err := s.Config.GetBackendList(c.Request.Context(), req)
		s.output(c, &resp, err)
	})

	router.GET("config/storage_list", func(c *gin.Context) {
		var req service.GetStorageListRequest
		if err := c.Bind(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		resp, err := s.Config.GetStorageList(c.Request.Context(), req)
		s.output(c, &resp, err)
	})

	router.POST("logger/search", func(c *gin.Context) {
		var req service.SearchRequest
		if err := c.Bind(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		resp, err := s.Logger.Search(c.Request.Context(), req)
		s.output(c, &resp, err)
	})
}

func (s *Application) output(c *gin.Context, resp interface{}, err error) {
	code, message := 0, "成功"
	if err != nil {
		code = 1
		message = err.Error()
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    resp,
	})
}
