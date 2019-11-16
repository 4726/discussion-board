package main

import (
	"fmt"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/4726/discussion-board/services/common"
)

const logInfoKey = "log info"

type RestAPI struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	api.engine = engine
	api.engine.Use(gin.Recovery())
	api.engine.Use(log.RequestMiddleware())
	api.setRoutes()
	common.AddMonitorHandler(api.engine)

	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&models.Comment{}, &models.Post{})
	api.db = db

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/post/:postid", func(ctx *gin.Context) {
		GetFullPost(a.db, ctx)
	})

	a.engine.GET("/posts", func(ctx *gin.Context) {
		GetPosts(a.db, ctx)
	})
}
func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
