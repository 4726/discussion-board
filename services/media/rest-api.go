package main

import (
	"github.com/gin-gonic/gin"
)

type RestAPI struct {
	mc     *MinioClient
	engine *gin.Engine
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}
	mc, err := NewMinioClient(cfg)
	if err != nil {
		return nil, err
	}
	api.mc = mc

	engine := gin.Default()
	api.engine = engine
	api.setRoutes()

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.POST("/upload", func(ctx *gin.Context) {
		Upload(a.mc, ctx)
	})

	a.engine.POST("/remove/:name", func(ctx *gin.Context) {
		Remove(a.mc, ctx)
	})

	a.engine.GET("/info", func(ctx *gin.Context) {
		Info(a.mc, ctx)
	})
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
