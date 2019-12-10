package main

import (
	"fmt"
	"net"

	"github.com/4726/discussion-board/services/likes/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

type GRPCApi struct {
	grpc     *grpc.Server
	handlers *GRPCHandlers //for testing
}

func NewGRPCApi(cfg Config) (*GRPCApi, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&CommentLike{}, &PostLike{})

	server := grpc.NewServer()
	handlers := &GRPCHandlers{db}
	pb.RegisterLikesServer(server, handlers)

	return &GRPCApi{server, handlers}, err
}

func (a *GRPCApi) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return a.grpc.Serve(lis)
}
