package main

import (
	"fmt"
	"net"

	"github.com/4726/discussion-board/services/posts/models"
	"github.com/4726/discussion-board/services/posts/read/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

type Api struct {
	db   *gorm.DB
	grpc *grpc.Server
}

func NewApi(cfg Config) (*Api, error) {
	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&models.Comment{}, &models.Post{})

	server := grpc.NewServer()
	handlers := &Handlers{db}
	pb.RegisterPostsReadServer(server, handlers)

	return &Api{db, server}, err
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return a.grpc.Serve(lis)
}
