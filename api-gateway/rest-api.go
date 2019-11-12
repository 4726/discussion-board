package main


import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RestAPI struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	api.setRoutes()

	return api, nil
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/post/:postid", func(ctx *gin.Context) {
		GetPost(ctx)
	})

	a.engine.GET("/posts", func(ctx *gin.Context) {
		GetPosts(ctx)
	})

	a.engine.POST("/post", func(ctx *gin.Context) {
		CreatePost(ctx)
	})

	a.engine.DELETE("/post/:postid", func(ctx *gin.Context) {
		DeletePost(ctx)
	})

	a.engine.POST("/post/like", func(ctx *gin.Context) {
		LikePost(ctx)
	})

	a.engine.POST("/post/unlike", func(ctx *gin.Context) {
		UnlikePost(ctx)
	})

	a.engine.POST("/comment", func(ctx *gin.Context) {
		AddComment(ctx)
	})

	a.engine.POST("/comment/like", func(ctx *gin.Context) {
		LikeComment(ctx)
	})

	a.engine.POST("/comment/clear", func(ctx *gin.Context) {
		ClearComment(ctx)
	})

	a.engine.POST("/search", func(ctx *gin.Context) {
		Search(ctx)
	})

	a.engine.GET("/register", func(ctx *gin.Context) {
		RegisterGET(ctx)
	})

	a.engine.POST("/register", func(ctx *gin.Context) {
		RegisterPOST(ctx)
	})

	a.engine.GET("/login", func(ctx *gin.Context) {
		LoginGET(ctx)
	})

	a.engine.POST("/login", func(ctx *gin.Context) {
		LoginPOST(ctx)
	})

	a.engine.POST("/changepassword", func(ctx *gin.Context) {
		ChangePassword(ctx)
	})

	a.engine.GET("/profile/:userid", func(ctx *gin.Context) {
		GetProfile(ctx)
	})

	a.engine.POST("/profile/update", func(ctx *gin.Context) {
		UpdateProfile(ctx)
	})
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
