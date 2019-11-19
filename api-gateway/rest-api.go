package main

import (
	"github.com/4726/discussion-board/services/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

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
	// api.engine.Use(log.RequestMiddleware())
	// api.setRoutes()
	api.setMockRoutes()
	common.AddMonitorHandler(api.engine)

	return api, nil
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/post/:postid", func(ctx *gin.Context) {
		GetPost(ctx)
	})

	a.engine.GET("/posts/:page", func(ctx *gin.Context) {
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

	a.engine.POST("/comment/unlike", func(ctx *gin.Context) {
		UnlikeComment(ctx)
	})

	a.engine.POST("/comment/clear", func(ctx *gin.Context) {
		ClearComment(ctx)
	})

	a.engine.GET("/search", func(ctx *gin.Context) {
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

func (a *RestAPI) setMockRoutes() {
	a.engine.GET("/post/:postid", func(ctx *gin.Context) {
		GetPostMock(ctx)
	})

	a.engine.GET("/posts/:page", func(ctx *gin.Context) {
		GetPostsMock(ctx)
	})

	a.engine.POST("/post", func(ctx *gin.Context) {
		CreatePostMock(ctx)
	})

	a.engine.DELETE("/post/:postid", func(ctx *gin.Context) {
		DeletePostMock(ctx)
	})

	a.engine.POST("/post/like", func(ctx *gin.Context) {
		LikePostMock(ctx)
	})

	a.engine.POST("/post/unlike", func(ctx *gin.Context) {
		UnlikePostMock(ctx)
	})

	a.engine.POST("/comment", func(ctx *gin.Context) {
		AddCommentMock(ctx)
	})

	a.engine.POST("/comment/like", func(ctx *gin.Context) {
		LikeCommentMock(ctx)
	})

	a.engine.POST("/comment/unlike", func(ctx *gin.Context) {
		UnlikeCommentMock(ctx)
	})

	a.engine.POST("/comment/clear", func(ctx *gin.Context) {
		ClearCommentMock(ctx)
	})

	a.engine.GET("/search", func(ctx *gin.Context) {
		SearchMock(ctx)
	})

	a.engine.GET("/register", func(ctx *gin.Context) {
		RegisterGETMock(ctx)
	})

	a.engine.POST("/register", func(ctx *gin.Context) {
		RegisterPOSTMock(ctx)
	})

	a.engine.GET("/login", func(ctx *gin.Context) {
		LoginGETMock(ctx)
	})

	a.engine.POST("/login", func(ctx *gin.Context) {
		LoginPOSTMock(ctx)
	})

	a.engine.POST("/changepassword", func(ctx *gin.Context) {
		ChangePasswordMock(ctx)
	})

	a.engine.GET("/profile/:userid", func(ctx *gin.Context) {
		GetProfileMock(ctx)
	})

	a.engine.POST("/profile/update", func(ctx *gin.Context) {
		UpdateProfileMock(ctx)
	})
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
