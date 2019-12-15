package main

import (
	"fmt"
	"net"

	"os"
	"os/signal"
	"syscall"

	"github.com/4726/discussion-board/services/common"
	pb "github.com/4726/discussion-board/services/likes/pb"
	_ "github.com/go-sql-driver/mysql"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Api struct {
	grpc     *grpc.Server
	handlers *Handlers //for testing
}

func NewApi(cfg Config) (*Api, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&CommentLike{}, &PostLike{})

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	shutdownCh := make(chan error, 1)
	go func() {
		sig := <-c
		a.grpc.GracefulStop()
		shutdownCh <- fmt.Errorf(sig.String())
	}()

	serveCh := make(chan error, 1)
	go func() {
		err := a.grpc.Serve(lis)
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		return err
	case err := <-shutdownCh:
		return err
	}
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
	creds, err := credentials.NewServerTLSFromFile(cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		return nil, err
	}

	return grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			common.IPWhiteListUnaryServerInterceptor(cfg.IPWhitelist),
			otgrpc.UnaryServerInterceptor(),
		)),
	), nil
}
