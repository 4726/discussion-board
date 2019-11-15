package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RestAPI struct {
	esc    *ESClient
	engine *gin.Engine
}

const logInfoKey = "log info"

func NewRestAPI(escIndexName, escAddr string) (*RestAPI, error) {
	gin.SetMode(gin.ReleaseMode)

	api := &RestAPI{}
	esc, err := NewESClient(escIndexName, escAddr)
	if err != nil {
		return nil, err
	}
	api.esc = esc

	engine := gin.New()
	api.engine = engine
	api.engine.Use(gin.Recovery())
	api.engine.Use(log.RequestMiddleware())
	api.setRoutes()
	api.setMonitorRoute()

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.POST("/index", func(ctx *gin.Context) {
		Index(a.esc, ctx)
	})

	a.engine.GET("/search", func(ctx *gin.Context) {
		Search(a.esc, ctx)
	})

	a.engine.POST("/update/likes", func(ctx *gin.Context) {
		UpdateLikes(a.esc, ctx)
	})

	a.engine.POST("/deletepost", func(ctx *gin.Context) {
		Delete(a.esc, ctx)
	})

	a.engine.POST("/update/lastupdate", func(ctx *gin.Context) {
		UpdateLastUpdate(a.esc, ctx)
	})
}

func (a *RestAPI) setMonitorRoute() {
	a.engine.Any("/metrics", gin.WrapH(promhttp.Handler()))
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
