package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type RestAPI struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	engine := gin.Default()
	api.engine = engine
	api.setRoutes()

	s := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Comment{}, &Post{})
	api.db = db

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.POST("/create/post", func(ctx *gin.Context) {
		CreatePost(a.db, ctx)
	})

	a.engine.POST("/delete/post", func(ctx *gin.Context) {
		DeletePost(a.db, ctx)
	})

	a.engine.POST("/update/post", func(ctx *gin.Context) {
		UpdatePostLikes(a.db, ctx)
	})

	a.engine.POST("/comment", func(ctx *gin.Context) {
		AddComment(a.db, ctx)
	})

	a.engine.POST("/comment/clear", func(ctx *gin.Context) {
		ClearComment(a.db, ctx)
	})

	a.engine.POST("/comment/likes", func(ctx *gin.Context) {
		UpdateCommentLikes(a.db, ctx)
	})
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
