package main

import (
	"fmt"
	"github.com/4726/discussion-board/services/posts/models"
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
	db.AutoMigrate(&models.Comment{}, &models.Post{})
	// deleting a post will also delete all of the post's comments
	db.Model(&models.Comment{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
	api.db = db

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.POST("/post/create", func(ctx *gin.Context) {
		CreatePost(a.db, ctx)
	})

	a.engine.POST("/post/delete", func(ctx *gin.Context) {
		DeletePost(a.db, ctx)
	})

	a.engine.POST("/post/likes", func(ctx *gin.Context) {
		UpdatePostLikes(a.db, ctx)
	})

	a.engine.POST("/comment/create", func(ctx *gin.Context) {
		CreateComment(a.db, ctx)
	})

	a.engine.POST("/comment/clear", func(ctx *gin.Context) {
		ClearComment(a.db, ctx)
	})

	a.engine.POST("/comment/likes", func(ctx *gin.Context) {
		UpdateCommentLikes(a.db, ctx)
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
