package main

import (
	"github.com/4726/discussion-board/api-gateway/pb/likes"
	postsread "github.com/4726/discussion-board/api-gateway/pb/posts-read"
	postswrite "github.com/4726/discussion-board/api-gateway/pb/posts-write"
	"github.com/4726/discussion-board/api-gateway/pb/search"
	"github.com/4726/discussion-board/api-gateway/pb/user"
	"github.com/4726/discussion-board/services/common"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type RestAPI struct {
	engine      *gin.Engine
	grpcClients GRPCClients
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	engine := gin.New()
	// gin.SetMode(gin.ReleaseMode)
	api.engine = engine
	api.engine.Use(corsMiddleware())
	api.engine.Use(gin.Recovery())
	api.engine.Use(log.RequestMiddleware())

	api.setupGRPCClients()

	// api.setRoutes()
	api.setMockRoutes()
	common.AddMonitorHandler(api.engine)

	return api, nil
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/post/:postid", func(ctx *gin.Context) {
		GetPost(ctx, a.grpcClients)
	})

	a.engine.GET("/posts", func(ctx *gin.Context) {
		GetPosts(ctx, a.grpcClients)
	})

	a.engine.POST("/post", func(ctx *gin.Context) {
		CreatePost(ctx, a.grpcClients)
	})

	a.engine.POST("/post/delete", func(ctx *gin.Context) {
		DeletePost(ctx, a.grpcClients)
	})

	a.engine.POST("/post/like", func(ctx *gin.Context) {
		LikePost(ctx, a.grpcClients)
	})

	a.engine.POST("/post/unlike", func(ctx *gin.Context) {
		UnlikePost(ctx, a.grpcClients)
	})

	a.engine.POST("/comment", func(ctx *gin.Context) {
		AddComment(ctx, a.grpcClients)
	})

	a.engine.POST("/comment/like", func(ctx *gin.Context) {
		LikeComment(ctx, a.grpcClients)
	})

	a.engine.POST("/comment/unlike", func(ctx *gin.Context) {
		UnlikeComment(ctx, a.grpcClients)
	})

	a.engine.POST("/comment/clear", func(ctx *gin.Context) {
		ClearComment(ctx, a.grpcClients)
	})

	a.engine.GET("/search", func(ctx *gin.Context) {
		Search(ctx, a.grpcClients)
	})

	a.engine.GET("/register", func(ctx *gin.Context) {
		RegisterGET(ctx)
	})

	a.engine.POST("/register", func(ctx *gin.Context) {
		RegisterPOST(ctx, a.grpcClients)
	})

	a.engine.GET("/login", func(ctx *gin.Context) {
		LoginGET(ctx)
	})

	a.engine.POST("/login", func(ctx *gin.Context) {
		LoginPOST(ctx, a.grpcClients)
	})

	a.engine.POST("/changepassword", func(ctx *gin.Context) {
		ChangePassword(ctx, a.grpcClients)
	})

	a.engine.GET("/profile/:userid", func(ctx *gin.Context) {
		GetProfile(ctx, a.grpcClients)
	})

	a.engine.POST("/profile/update", func(ctx *gin.Context) {
		UpdateProfile(ctx, a.grpcClients)
	})

	a.engine.GET("/userid", func(ctx *gin.Context) {
		// UserID(ctx)
	})
}

func (a *RestAPI) setMockRoutes() {
	a.engine.GET("/post/:postid", func(ctx *gin.Context) {
		GetPostMock(ctx)
	})

	a.engine.GET("/posts", func(ctx *gin.Context) {
		GetPostsMock(ctx)
	})

	a.engine.POST("/post", func(ctx *gin.Context) {
		CreatePostMock(ctx)
	})

	a.engine.POST("/post/delete", func(ctx *gin.Context) {
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

	a.engine.GET("/userid", func(ctx *gin.Context) {
		UserIDMock(ctx)
	})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (a *RestAPI) setupGRPCClients() {
	userConn, err := grpc.Dial(UserServiceAddr(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient := user.NewUserClient(userConn)
	searchConn, err := grpc.Dial(SearchServiceAddr(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	searchClient := search.NewSearchClient(searchConn)
	likesConn, err := grpc.Dial(LikesServiceAddr(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	likesClient := likes.NewLikesClient(likesConn)
	postreadConn, err := grpc.Dial(PostsReadServiceAddr(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	postsreadClient := postsread.NewPostsReadClient(postreadConn)
	postWriteConn, err := grpc.Dial(PostsWriteServiceAddr(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	postswriteClient := postswrite.NewPostsWriteClient(postWriteConn)
	a.grpcClients = GRPCClients{
		searchClient,
		userClient,
		likesClient,
		postsreadClient,
		postswriteClient,
	}
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}
