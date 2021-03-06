package main

import (
	"fmt"
	"net"

	"github.com/4726/discussion-board/services/common"
	pb "github.com/4726/discussion-board/services/likes/pb"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type Api struct {
	grpc     *grpc.Server
	handlers *Handlers //for testing
}

func NewApi(cfg Config) (*Api, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	log.Entry().Infof("connecting to database: %s", s)
	db, err := common.OpenDB("mysql", s)
	if err != nil {
		return nil, err
	}
	log.Entry().Infof("successfully connected to database: %s", s)
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

	log.Entry().Infof("server running on addr: %s", addr)

	return common.RunGRPCWithGracefulShutdown(a.grpc, lis)
}

func tlsGRPC(cfg Config) (*grpc.Server, error) {
	opts := common.GRPCOptions{cfg.IPWhitelist, cfg.TLSCert, cfg.TLSKey, log.Entry()}
	return common.DefaultGRPCServer(opts)
}
