package main

import (
	"fmt"
	"net"

	pb "github.com/4726/discussion-board/services/likes/pb"
	_ "github.com/go-sql-driver/mysql"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

type Api struct {
	grpc     *grpc.Server
	handlers *Handlers //for testing
}

func NewApi(cfg Config) (*Api, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&CommentLike{}, &PostLike{})

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.UnaryServerInterceptor()))
	handlers := &Handlers{db}
	pb.RegisterLikesServer(server, handlers)

	return &Api{server, handlers}, err
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return a.grpc.Serve(lis)
}
