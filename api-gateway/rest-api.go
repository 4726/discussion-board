package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	likes "github.com/4726/discussion-board/api-gateway/pb/likes"
	postsread "github.com/4726/discussion-board/api-gateway/pb/posts-read"
	postswrite "github.com/4726/discussion-board/api-gateway/pb/posts-write"
	search "github.com/4726/discussion-board/api-gateway/pb/search"
	user "github.com/4726/discussion-board/api-gateway/pb/user"
	"github.com/4726/discussion-board/services/common"
	"github.com/gin-gonic/gin"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type RestAPI struct {
	engine      *gin.Engine
	grpcClients GRPCClients
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	engine := gin.New()
	gin.SetMode(gin.ReleaseMode)
	api.engine = engine
	api.engine.Use(corsMiddleware())
	api.engine.Use(gin.Recovery())
	// api.engine.Use(log.RequestMiddleware())

	api.setupGRPCClients(cfg)
	api.setRoutes()

	// api.setMockRoutes()
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
		UserIdGET(ctx)
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
		UserIdGETMock(ctx)
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

func (a *RestAPI) setupGRPCClients(cfg Config) {
	creds, _ := credentials.NewClientTLSFromFile(cfg.UserService.TLSCert, cfg.UserService.TLSServerName)
	userConn, _ := grpc.Dial(
		cfg.UserService.Addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otgrpc.UnaryClientInterceptor()),
	)
	userClient := user.NewUserClient(userConn)

	creds, _ = credentials.NewClientTLSFromFile(cfg.SearchService.TLSCert, cfg.SearchService.TLSServerName)
	searchConn, _ := grpc.Dial(
		cfg.SearchService.Addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otgrpc.UnaryClientInterceptor()),
	)
	searchClient := search.NewSearchClient(searchConn)

	creds, _ = credentials.NewClientTLSFromFile(cfg.LikesService.TLSCert, cfg.LikesService.TLSServerName)
	likesConn, _ := grpc.Dial(
		cfg.LikesService.Addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otgrpc.UnaryClientInterceptor()),
	)
	likesClient := likes.NewLikesClient(likesConn)

	creds, _ = credentials.NewClientTLSFromFile(cfg.PostsReadService.TLSCert, cfg.PostsReadService.TLSServerName)
	postreadConn, _ := grpc.Dial(
		cfg.PostsReadService.Addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otgrpc.UnaryClientInterceptor()),
	)
	postsreadClient := postsread.NewPostsReadClient(postreadConn)

	creds, _ = credentials.NewClientTLSFromFile(cfg.PostsWriteService.TLSCert, cfg.PostsWriteService.TLSServerName)
	postWriteConn, _ := grpc.Dial(
		cfg.PostsWriteService.Addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otgrpc.UnaryClientInterceptor()),
	)
	postswriteClient := postswrite.NewPostsWriteClient(postWriteConn)

	a.grpcClients = GRPCClients{
		searchClient,
		userClient,
		likesClient,
		postsreadClient,
		postswriteClient,
	}
}

func (a *RestAPI) Run(addr, tlsCert, tlsKey string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: a.engine,
	}

	log.Entry().Infof("server running on addr: %s", addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	shutdownCh := make(chan error, 1)
	go func() {
		sig := <-c
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			shutdownCh <- err
		} else {
			shutdownCh <- fmt.Errorf(sig.String())
		}
	}()

	serveCh := make(chan error, 1)
	go func() {
		err := s.ListenAndServeTLS(tlsCert, tlsKey)
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		return err
	case err := <-shutdownCh:
		return err
	}

}
