package main

import (
	"github.com/gin-gonic/gin"
)

type RestAPI struct {
	esc    *ESClient
	engine *gin.Engine
}

func NewRestAPI() (*RestAPI, error) {
	api := &RestAPI{}
	esc, err := NewESClient()
	if err != nil {
		return nil, err
	}
	api.esc = esc

	engine := gin.Default()
	api.engine = engine
	api.setRoutes()

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/search", func(ctx *gin.Context) {
		Search(a.esc, ctx)
	})

	a.engine.POST("/index", func(ctx *gin.Context) {
		Index(a.esc, ctx)
	})

	a.engine.POST("/updatelikes", func(ctx *gin.Context) {
		UpdateLikes(a.esc, ctx)
	})

	a.engine.DELETE("/index/:id", func(ctx *gin.Context) {
		Delete(a.esc, ctx)
	})
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
