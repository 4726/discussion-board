package main

import (
	"fmt"
	"github.com/4726/discussion-board/posts/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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

	s := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.DBName)

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
			standardRequestLoggingEntry(c).Error(logMessage)
			return
		}

		if c.Writer.Status() == http.StatusOK {
			standardRequestLoggingEntry(c).Info(logMessage)
			return
		}

		if c.Writer.Status() == http.StatusNotFound {
			standardRequestLoggingEntry(c).Info(logMessage)
			return
		}

		if c.Writer.Status() == http.StatusBadRequest {
			standardRequestLoggingEntry(c).Warn(logMessage)
			return
		}
	}
}
