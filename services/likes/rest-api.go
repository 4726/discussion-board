package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const logInfoKey = "log info"

type RestAPI struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()
	api.engine = engine
	api.engine.Use(api.logRequestsMiddleware())
	api.setRoutes()
	api.setMonitorRoute()

	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&CommentLike{}, &PostLike{})
	api.db = db

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.POST("/post/like", func(ctx *gin.Context) {
		LikePost(a.db, ctx)
	})

	a.engine.POST("/post/unlike", func(ctx *gin.Context) {
		UnlikePost(a.db, ctx)
	})

	a.engine.POST("/comment/like", func(ctx *gin.Context) {
		LikeComment(a.db, ctx)
	})

	a.engine.POST("/comment/unlike", func(ctx *gin.Context) {
		UnlikeComment(a.db, ctx)
	})

	a.engine.GET("/post/likes", func(ctx *gin.Context) {
		GetMultiplePostLikes(a.db, ctx)
	})

	a.engine.GET("/comment/likes", func(ctx *gin.Context) {
		GetMultipleCommentLikes(a.db, ctx)
	})
}

func (a *RestAPI) setMonitorRoute() {
	a.engine.Any("/metrics", gin.WrapH(promhttp.Handler()))
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}

func (a *RestAPI) logRequestsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logMessage := ""
		i, ok := c.Get(logInfoKey)
		if ok {
			switch v := i.(type) {
			case string:
				logMessage = v
			case error:
				logMessage = v.Error()
			default:
			}
		}

		if c.Writer.Status() == http.StatusInternalServerError {
			log.HTTPRequestEntry(c).Error(logMessage)
			return
		}

		if c.Writer.Status() == http.StatusOK {
			log.HTTPRequestEntry(c).Info(logMessage)
			return
		}

		if c.Writer.Status() == http.StatusBadRequest {
			log.HTTPRequestEntry(c).Warn(logMessage)
			return
		}
	}
}