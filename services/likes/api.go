package main

import (
	"fmt"
	"net"

	"github.com/4726/discussion-board/services/common"
	pb "github.com/4726/discussion-board/services/likes/pb"
	"github.com/cenkalti/backoff/v3"
	_ "github.com/go-sql-driver/mysql"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

type Api struct {
	grpc     *grpc.Server
	handlers *Handlers //for testing
}

func NewApi(cfg Config) (*Api, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	var db *gorm.DB
	op := func() error {
		var err error
		db, err = gorm.Open("mysql", s)
		if err != nil {
			return err
		}
		// db.LogMode(true)
		db.AutoMigrate(&CommentLike{}, &PostLike{})
		return nil
	}

	err := backoff.Retry(op, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	server, err := tlsGRPC(cfg)
	if err != nil {
		return nil, err
	}
	handlers := &Handlers{db}
	pb.RegisterLikesServer(server, handlers)

	return &Api{server, handlers}, err
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Entry().Infof("server running on addr: %s", addr)

	return common.RunGRPCWithGracefulShutdown(a.grpc, lis)
}

func tcpGRPC(cfg Config) (*grpc.Server, error) {
	return grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			common.IPWhiteListUnaryServerInterceptor(cfg.IPWhitelist),
			otgrpc.UnaryServerInterceptor(),
		)),
	), nil
}

func tlsGRPC(cfg Config) (*grpc.Server, error) {
	opts := common.GRPCOptions{cfg.IPWhitelist, cfg.TLSCert, cfg.TLSKey}
	return common.DefaultGRPCServer(opts)
}
