package main

import (
	"fmt"
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

	engine := gin.New()
	gin.SetMode(gin.ReleaseMode)
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
	db.AutoMigrate(&Auth{}, &Profile{})
	api.db = db

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/profile/:userid", func(ctx *gin.Context) {
		GetProfile(a.db, ctx)
	})

	a.engine.POST("/login", func(ctx *gin.Context) {
		ValidLogin(a.db, ctx)
	})

	a.engine.POST("/account", func(ctx *gin.Context) {
		CreateAccount(a.db, ctx)
	})

	a.engine.POST("/profile/update", func(ctx *gin.Context) {
		UpdateProfile(a.db, ctx)
	})

	a.engine.POST("/password", func(ctx *gin.Context) {
		ChangePassword(a.db, ctx)
	})
}

func (a *RestAPI) Run(addr string) error {
	log.Entry().Info("starting service on addr: " + addr)
	return a.engine.Run(addr)
}
